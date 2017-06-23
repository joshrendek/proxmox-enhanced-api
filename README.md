# proxmox-enhanced-api

`proxmox-enhanced-api` will consist of the following components in a single binary. It can be deployed on the server, ran once to install systemd template, and will run itself as a daemon from there.

Config file will live most likely under /etc/proxmox-enhanced-api/

# Authentication

It will use the same auth mechanism as the Proxmox API and will just proxy auth requests to verify them against the regular API. It will also provide the ability to just use the user/pass instead of a ticket.

## Guest VM IP API

- [ ] Ability to list all VMs on host with MACs
- [ ] Ability to query a VM by name or MAC to get the IP of the machine

## VM to DNS Registration

- [ ] Take all VMs and provide the ability to register them to a DNS Provider configured in the `config file`
- [ ] Initial providers:
  - [ ] CloudFlare
  - [ ] Route53
