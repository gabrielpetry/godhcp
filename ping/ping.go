package ping

// based on https://gist.github.com/lmas/c13d1c9de3b2224f9c26435eb56e6ef3

import (
	"fmt"
	Config "go-dhcpdump/config"
	"go-dhcpdump/log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Response struct {
	Addr *net.IPAddr
	Rtt  time.Duration
}

const (
	// can't import "internal" packages
	ProtocolICMP = 1
	//ProtocolIPv6ICMP = 58
)

var config = Config.GetInstance()

const waitTimeout = time.Second * 2
const errorLimit = 3

func PingUntilFail(target string) (int, error) {
	errorCounter := 0
	var err error

	for {
		if errorCounter >= errorLimit {
			log.Debug("target stopped responding", target, errorCounter)
			break
		}
		time.Sleep(500 * time.Millisecond)
		// @TODO: Retry with the Echo request after fails timestamp
		err = SendIcmpTimestamp(target)
		// server respond without error
		if err == nil {
			errorCounter = 0
			continue
		}
		errorCounter += 1
	}

	return errorCounter, err
}

func SendIcmpTimestamp(target string) error {
	buf := make([]byte, 16)
	m := icmp.Message{
		Type: ipv4.ICMPTypeTimestamp,
		Code: 0,
		Body: &icmp.DefaultMessageBody{Data: buf},
	}

	b, err := m.Marshal(nil)
	if err != nil {
		return fmt.Errorf("Error marshaling message: %v", err)
	}

	return sendPing(b, target)
}

func SendIcmpEcho(target string) error {
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, //<< uint(seq), // TODO
			Data: []byte(""),
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return fmt.Errorf("Error marshaling message: %v", err)
	}

	return sendPing(b, target)
}

func sendPing(b []byte, target string) error {
	// Start listening for icmp replies
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("Error creating listenPacket: %v", err)
	}
	defer c.Close()

	// Resolve any DNS (if used) and get the real IP of the target
	dst, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		return fmt.Errorf("Error resolving ip: %v", err)
	}

	// Send it
	n, err := c.WriteTo(b, dst)
	if err != nil {
		return fmt.Errorf("Error writing to socket: %v", err)
	} else if n != len(b) {
		return fmt.Errorf("Wrong Response: got %v; want %v", n, len(b))
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(waitTimeout))
	if err != nil {
		return fmt.Errorf("Error setting read deadline: %v", err)
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return fmt.Errorf("Error reading from socket: %v", err)
	}

	// Pack it up boys, we're done here
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return fmt.Errorf("Error parsing icmp message %v", err)
	}

	switch rm.Type {
	case ipv4.ICMPTypeTimestampReply:
	case ipv4.ICMPTypeEchoReply:
		return nil
	default:
		return fmt.Errorf("Wrong Response: got %v; want %v", rm, peer)
	}

	return nil
}
