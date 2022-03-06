package main

import (
	"go-dhcpdump/api"
	Config "go-dhcpdump/config"
	"go-dhcpdump/dhcpDump"
	"go-dhcpdump/dhcpMessage"
	"go-dhcpdump/log"
	"go-dhcpdump/ping"
	"time"
)

func main() {
	log.Info("Starting go-dhcpdump")

	config := Config.GetInstance()
	log.Info("Using config", config)

	go api.StartApiServer()

	go func() {
		time.Sleep(time.Second * 2)
		dhcp := dhcpMessage.DhcpdumpMessage{}
		onlineDevices := dhcp.GetOnlineDevices()

		log.Debug("Pingin old devices: ", onlineDevices)
		for _, v := range onlineDevices {
			err := ping.Ping(v.ClientIpAddress)
			if err == nil {
				v.Save()
				ping.AddJob(v)
				continue
			}
			dhcp.DisconnectDevice()
		}
	}()
	go dhcpDump.RunDhcpdump()

	ping.StartWorker()
}
