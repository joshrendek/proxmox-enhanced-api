# proxmox-enhanced-api

# Requirements

``` shell
apt-get update && apt-get install -y arp-scan
```

`proxmox-enhanced-api` will consist of the following components in a single binary. It can be deployed on the server, ran once to install systemd template, and will run itself as a daemon from there.

# Installation

Download the binary to `/root/` and run `./proxmox-enhanced-api -init` -- this will generate a config and systemd unit files. You can then edit the config file located in `/etc/proxmox-enhanced-api/config.toml` and restart the service using `systemctl restart proxmox-enhanced-api.service`.

*Example Configuration*:

``` bash
[proxmox]
host = "192.168.1.2"
node = "pve"
user = "root@pam"
pass = "foobar123"

[api]
# if skip_auth is true, you need to enter user/pass credentials under proxmox,
# otherwise you need to pass the user/pass a basic auth
skip_auth = true
```

# Accessing the service

You can access it by going to `http://YOUR_VM_IP:8080/vm`

Example with `skip_auth = true`:

``` bash
curl http://192.168.1.2:8080/vm
```

Example with basic auth:

``` bash
curl -u "root@pam":root http://192.168.1.2:8080/vm
```

# Example response

``` json
[
  {
    "name": "jenkins-slave",
    "vmid": 101,
    "status": "running",
    "mac_address": "4A:BB:E4:82:C9:4C",
    "ip_address": "192.168.1.183"
  },
  {
    "name": "testing",
    "vmid": 104,
    "status": "running",
    "mac_address": "BE:D9:B5:7A:42:19",
    "ip_address": "192.168.1.195"
  },
  {
    "name": "zoneminder",
    "vmid": 103,
    "status": "running",
    "mac_address": "3A:11:9E:32:9B:D4",
    "ip_address": "192.168.1.167"
  }
]
```

## Authentication

It will use the same auth mechanism as the Proxmox API and will just proxy auth requests to verify them against the regular API. It will also provide the ability to just use the user/pass instead of a ticket.

## Guest VM IP API

- [x] Ability to list all VMs on host with MACs and IPs

## VM to DNS Registration

- [ ] Take all VMs and provide the ability to register them to a DNS Provider configured in the `config file`
- [ ] Initial providers:
  - [ ] CloudFlare
  - [ ] Route53
