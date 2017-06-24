package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Proxmox struct {
	Username string
	Password string
	Host     string
	Ticket   string
	CSRF     string
}

var (
	arpLock   sync.Mutex
	transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	macList = map[string]string{}
)

func StartArper() {
	for {
		fmt.Println("[arp] refresh started")
		arpLock.Lock()
		macList = MacList()
		arpLock.Unlock()
		fmt.Println("[arp] refresh finished")
		time.Sleep(30 * time.Second)
	}
}

func NewProxmox(user, pass, host string) (*Proxmox, error) {

	// get a ticket

	p := &Proxmox{Username: user, Password: pass, Host: host}

	form := url.Values{
		"username": {p.Username},
		"password": {p.Password},
	}

	macList = MacList()

	target := fmt.Sprintf("https://%s:8006/api2/json/access/ticket", host)
	req, err := http.NewRequest("POST", target, bytes.NewBufferString(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	ticket := Ticket{}
	err = json.Unmarshal(response, &ticket)
	p.Ticket = ticket.Data.Ticket
	p.CSRF = ticket.Data.CSRFPreventionToken

	return p, nil
}

func (p *Proxmox) VirtualMachines() ([]VirtualMachine, error) {
	//fmt.Printf("MACLIST: %+v\n", macList)
	qemus, err := p.Qemu()
	if err != nil {
		return nil, err
	}
	ret := []VirtualMachine{}
	for _, q := range qemus {
		if q.Status == "stopped" {
			continue
		}
		tmp := VirtualMachine{Name: q.Name, Vmid: q.Vmid, Status: q.Status}
		config, _ := p.Config(q)
		m := strings.Split(config.Net0, ",")
		if len(m) == 0 {
			continue
		}
		m2 := strings.Split(m[0], "=")
		if len(m2) == 0 {
			continue
		}
		tmp.MacAddress = m2[1]
		tmp.IPAddress = macList[strings.ToLower(tmp.MacAddress)]
		ret = append(ret, tmp)
	}

	return ret, nil
}

func MacList() map[string]string {
	ret := map[string]string{}
	cmd := exec.Command("arp-scan", "--interface=vmbr0", "--localnet")
	out, err := cmd.CombinedOutput()
	//fmt.Println("OUT: ", string(out))
	if err != nil {
		panic(err)
	}
	s := strings.Split(string(out), "\n")
	for _, m := range s {
		tmp := strings.Fields(m)
		if len(tmp) == 3 {
			ret[tmp[1]] = tmp[0]
		}
	}
	return ret
}

func (p *Proxmox) Config(q Qemu) (*Config, error) {
	target := fmt.Sprintf("https://%s:8006/api2/json/nodes/pve/qemu/%d/config", p.Host, q.Vmid)
	req, err := http.NewRequest("GET", target, nil)
	req.Header.Add("Cookie", fmt.Sprintf("PVEAuthCookie=%s", p.Ticket))
	req.Header.Add("CSRFPreventionToken", p.CSRF)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	configData := ConfigData{}
	if err := json.Unmarshal(response, &configData); err != nil {
		return nil, err
	}
	return &configData.Config, nil
}

func (p *Proxmox) Qemu() ([]Qemu, error) {
	target := fmt.Sprintf("https://%s:8006/api2/json/nodes/pve/qemu", p.Host)
	req, err := http.NewRequest("GET", target, nil)
	req.Header.Add("Cookie", fmt.Sprintf("PVEAuthCookie=%s", p.Ticket))
	req.Header.Add("CSRFPreventionToken", p.CSRF)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	qData := QemuData{}
	if err := json.Unmarshal(response, &qData); err != nil {
		return nil, err
	}

	return qData.Qemus, nil
}
