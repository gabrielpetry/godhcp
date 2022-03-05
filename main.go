package main

import (
	Config "go-dhcpdump/config"
	dhcpDump "go-dhcpdump/dhcpdump"
	"go-dhcpdump/log"
)

func main() {
	log.Info("Starting go-dhcpdump")

	config := Config.GetInstance()
	log.Info("Using config", config)

	// go api.StartApiServer()

	dhcp := dhcpDump.DhcpdumpMessage{}
	onlineDevices := dhcp.GetOnlineDevices()
	func() {
		for _, v := range onlineDevices {
			v.SinglePing()
		}
	}()

	dhcpDump.RunDhcpdump()
}
