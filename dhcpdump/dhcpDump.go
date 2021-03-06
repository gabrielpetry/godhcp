package dhcpDump

import (
	Config "go-dhcpdump/config"
	"go-dhcpdump/dhcpMessage"
	"go-dhcpdump/log"
	"go-dhcpdump/ping"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var config = Config.GetInstance()

func RunDhcpdump() {
	iface := Config.GetInstance().Dhcpdump.Interface
	handle, err := pcap.OpenLive(iface, 1500, true, 5*time.Second)

	if err != nil { // optional
		panic(err)
	}
	// dhcpdump filter
	handle.SetBPFFilter("udp and (port bootpc or port bootps)")
	// handle packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handleDhcpPacket(packet) // Do something with a packet here.
	}
}

func handleDhcpPacket(message gopacket.Packet) {
	dhcp := dhcpMessage.DhcpdumpMessage{}
	dhcp.Parse(message.String())
	err := dhcp.FilterDhcpPacket()
	if err != nil {
		log.Error("Device is not filtered: ", err)
		return
	}

	dhcp.Save()
	ping.AddJob(dhcp)
}
