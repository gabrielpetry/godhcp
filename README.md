# Godhcp


This is bizzarre project

## The idea

My ISP provided a router without access to way to make a static ip for each device nor and api to consult based on hostname

This program is just this:

A dhcpDump ~kinda~ watcher, using libpcap

A pinger, using golangs icmp (trying both Echo and Timestamps request)

A queue (fast solution for pinging multiple devices, since i'm lazy to figure out a better approch)

Mqtt broker to send the status of the devices to a Hass server

## How to use

How to use on raspian

build on your machine with docker:
`make build_raspi`

install pcap on raspian
`sudo apt update && sudo apt install libpcap0.8 -y`

Supervisor conf
```conf
[program:godhcpdump]
command=/home/pi/godhcpdump
```

`sudo supervisorctl update ; sudo supervisorctl restart godhcpdump`

# TODOS:
[ ] improve doc
[ x ] make this program work
[ ] improve this program