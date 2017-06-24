package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

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

func main() {
	initBool := false
	flag.BoolVar(&initBool, "init", false, "init systemd / config")
	flag.Parse()

	viper.SetConfigName("config")                     // name of config file (without extension)
	viper.AddConfigPath("/etc/proxmox-enhanced-api/") // path to look for the config file in
	err := viper.ReadInConfig()                       // Find and read the config file
	if err != nil {                                   // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("NODE: ", viper.GetString("proxmox.node"))

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
		os.Exit(0)
	}
	// systemctl daemon-reload
	// systemctl stop proxmox-enhanced-api.service
	// systemctl start proxmox-enhanced-api.service

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
	r.Run()

}
