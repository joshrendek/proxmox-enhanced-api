package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/crackcomm/cloudflare"
	"github.com/gin-gonic/gin"
	"github.com/joshrendek/proxmox-enhanced-api/proxmox"
	"github.com/spf13/viper"
)

const systemdUnit = `
[Unit]
Description=proxmox-enhanced-api
After=network.target

[Service]
ExecStart=/root/proxmox-enhanced-api
Type=simple

[Install]
WantedBy=default.target
`

const peaConfig = `
[proxmox]
host = "192.168.1.2"
node = "pve"
user = "root@pam"
pass = "foobar123"

[api]
# if skip_auth is true, you need to enter user/pass credentials under proxmox,
# otherwise you need to pass the user/pass a basic auth
skip_auth = true

[dns]
zone = "example.com"

[cloudflare]
api_key = "123"
email = "me@example.com"
`

func main() {
	initBool := false
	flag.BoolVar(&initBool, "init", false, "init systemd / config")
	flag.Parse()

	if initBool {
		fmt.Println("[systemd] Writing unit")
		exec.Command("systemctl", "stop", "proxmox-enhanced-api.service").Run()
		err := ioutil.WriteFile("/etc/systemd/system/proxmox-enhanced-api.service", []byte(systemdUnit), 0644)
		if err != nil {
			log.Fatalln(err)
		}
		exec.Command("systemctl", "daemon-reload").Run()
		exec.Command("systemctl", "start", "proxmox-enhanced-api.service").Run()
		fmt.Println("[systemd] Finished writing unit")

		fmt.Println("[config] Writing config")
		err = os.Mkdir("/etc/proxmox-enhanced-api", os.FileMode(0644))
		if err != nil {
			log.Fatalln(err)
		}
		// init config
		if _, err := os.Stat("/etc/proxmox-enhanced-api/config.toml"); os.IsNotExist(err) {
			// path/to/whatever does not exist
			err := ioutil.WriteFile("/etc/proxmox-enhanced-api/config.toml", []byte(peaConfig), 0644)
			if err != nil {
				log.Fatalln(err)
			}
		}

		fmt.Println("[config] Finished writing config")

		os.Exit(0)
	}

	viper.SetConfigName("config")                     // name of config file (without extension)
	viper.AddConfigPath("/etc/proxmox-enhanced-api/") // path to look for the config file in
	err := viper.ReadInConfig()                       // Find and read the config file
	if err != nil {                                   // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("NODE: ", viper.GetString("proxmox.node"))
	// systemctl daemon-reload
	// systemctl stop proxmox-enhanced-api.service
	// systemctl start proxmox-enhanced-api.service
	go registerDNS()
	go proxmox.StartArper()
	r := gin.Default()
	r.GET("/vm", func(c *gin.Context) {
		vms := []proxmox.VirtualMachine{}
		var proxUser, proxPass string
		if viper.GetBool("api.skip_auth") {
			proxUser = viper.GetString("proxmox.user")
			proxPass = viper.GetString("proxmox.pass")
		} else {
			user, pass, _ := c.Request.BasicAuth()
			if len(user) == 0 || len(pass) == 0 {
				c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
				return
			}
			proxUser = user
			proxPass = pass
		}

		prox, err := proxmox.NewProxmox(proxUser, proxPass, viper.GetString("proxmox.host"))
		if err != nil {
			panic(err)
		}
		vms, err = prox.VirtualMachines()
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"status": "error with credentials"})
			return
		}
		c.JSON(http.StatusOK, vms)
	})
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	err = http.ListenAndServeTLS(fmt.Sprintf(":%s", port), "/etc/pve/nodes/pve/pve-ssl.pem", "/etc/pve/nodes/pve/pve-ssl.key", r)
	if err != nil {
		log.Fatal(err)
	}
}

func registerDNS() {
	localList := []proxmox.VirtualMachine{}
	proxUser := viper.GetString("proxmox.user")
	proxPass := viper.GetString("proxmox.pass")
	prox, err := proxmox.NewProxmox(proxUser, proxPass, viper.GetString("proxmox.host"))
	if err != nil {
		panic(err)
	}

	client := cloudflare.New(&cloudflare.Options{
		Email: viper.GetString("cloudflare.email"),
		Key:   viper.GetString("cloudflare.api_key"),
	})

	var zone *cloudflare.Zone

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	zones, err := client.Zones.List(ctx)
	if err != nil {
		log.Fatal(err)
	} else if len(zones) == 0 {
		log.Fatal("No zones were found")
	}

	for _, z := range zones {
		if z.Name == viper.GetString("dns.zone") {
			zone = z
		}
	}

	records, err := client.Records.List(ctx, zone.ID)
	for _, record := range records {
		log.Printf("RECORD: %#v", record)
	}

	for {
		localList, _ = prox.VirtualMachines()
		for _, v := range localList {
			fqdn := fmt.Sprintf("%s.%s", v.Name, viper.GetString("dns.zone"))
			log.Printf("DNS: %s -> %s", v.Name, fqdn)
			createRecord(fqdn, v.IPAddress, zone, client)
		}
		time.Sleep(30 * time.Second)
	}
}

func createRecord(name, ip string, zone *cloudflare.Zone, client *cloudflare.Client) {
	cfRecord := &cloudflare.Record{Type: "A", Name: name,
		Content: ip, ZoneID: zone.ID, ZoneName: zone.Name}
	log.Println("[dns] createRecord", name)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	err := client.Records.Create(ctx, cfRecord)
	if err.Error() == "The record already exists." {
		client.Records.Patch(ctx, cfRecord)
		log.Println("[dns] patching record", name)
	}

}
