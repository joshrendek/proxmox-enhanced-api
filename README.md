# proxmox-enhanced-api

# Requirements

``` shell
apt-get update && apt-get install -y arp-scan
```

`proxmox-enhanced-api` will consist of the following components in a single binary. It can be deployed on the server, ran once to install systemd template, and will run itself as a daemon from there.

# Security

Auth is currently done with user/pass or skipped (and user/pass need to be set in config, see below).

The server listens on port `8080` using the same SSL certification that proxmox uses for the UI, located at:

``` json
/etc/pve/nodes/pve/pve-ssl.pem
/etc/pve/nodes/pve/pve-ssl.key
```

Credentials are verified against the API over SSL as well.

Ticket auth is not supported currently.

# Installation

Download the binary to `/root/` and run `./proxmox-enhanced-api -init`.

This will generate a config and systemd unit files.
You can then edit the config file located in `/etc/proxmox-enhanced-api/config.toml` and restart the service using
`systemctl restart proxmox-enhanced-api.service`.

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

[dns]
zone = "example.com"

[cloudflare]
api_key = "123"
email = "me@example.com"
```

# Accessing the service

You can access it by going to `http://YOUR_VM_IP:8080/vm`

Example with `skip_auth = true`:

``` bash
curl -k https://192.168.1.2:8080/vm
```

Example with basic auth:

``` bash
curl -k -u "root@pam":root https://192.168.1.2:8080/vm
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

## Guest VM IP API

- [x] Ability to list all VMs on host with MACs and IPs

## VM to DNS Registration

- [x] Take all VMs and provide the ability to register them to a DNS Provider configured in the `config file`
- [ ] Initial providers:
  - [x] CloudFlare
  - [ ] Route53
