package ping

// import (
// 	"fmt"
// 	Config "go-dhcpdump/config"
// 	"go-dhcpdump/log"
// 	"net"
// 	"time"

// 	"golang.org/x/net/icmp"
// 	"golang.org/x/net/ipv4"
// )

// type Response struct {
// 	Addr *net.IPAddr
// 	Rtt  time.Duration
// }

// var config = Config.GetInstance()

// const waitTimeout = time.Second * 2
// const errorLimit = 3

// func Ping(target string) (int, error) {
// 	errorCounter := 0
// 	var err error

// 	for {
// 		if errorCounter >= errorLimit {
// 			log.Debug("target stopped responding", target, errorCounter)
// 			break
// 		}
// 		time.Sleep(500 * time.Millisecond)
// 		err = sendPing(target)
// 		// server respond without error
// 		if err == nil {
// 			log.Debug("Server responded ok")
// 			errorCounter = 0
// 			continue
// 		}
// 		log.Debug("Server responded with error", err)
// 		errorCounter += 1
// 	}

// 	return errorCounter, err
// }

// func sendPing(target string) error {
// 	// create the ip
// 	log.Debug("Sending one ping", target)
// 	ip, err := net.ResolveIPAddr("ip4", target)
// 	if err != nil {
// 		return fmt.Errorf("Ip unreachable %v", err)
// 	}

// 	// create the listener
// 	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
// 	defer conn.Close()

// 	// define the icmp timestamp body
// 	buf := make([]byte, 16)
// 	msg := icmp.Message{
// 		Type: ipv4.ICMPTypeTimestamp,
// 		Code: 0,
// 		Body: &icmp.DefaultMessageBody{Data: buf},
// 	}
// 	msg_bytes, err := msg.Marshal(nil)
// 	if err != nil {
// 		return fmt.Errorf("Error marshaling msg bytes %v", err)
// 	}

// 	// Write the message to the listening connection
// 	dest := &net.UDPAddr{IP: net.ParseIP(ip.String())}
// 	fmt.Printf("%v", *dest)
// 	log.Debug(*dest)
// 	if _, err := conn.WriteTo(msg_bytes, dest); err != nil {
// 		return fmt.Errorf("Error on WriteTo %v : %v", err, msg_bytes)
// 	}
// 	// return fmt.Errorf("Error on WriteTo %v : %v", err, msg_bytes)
// 	fmt.Printf("%v", msg_bytes)
// 	// set the timeout
// 	err = conn.SetReadDeadline(time.Now().Add(waitTimeout))
// 	if err != nil {
// 		return fmt.Errorf("Error on SetReadDeadline %v", err)
// 	}

// 	// get the reply
// 	reply := make([]byte, 1500)
// 	n, _, err := conn.ReadFrom(reply)
// 	if err != nil {
// 		return fmt.Errorf("Error on ReadFrom %v", err)
// 	}
// 	parsed_reply, err := icmp.ParseMessage(1, reply[:n])
// 	if err != nil {
// 		return fmt.Errorf("Error on ParseMessage %v", err)
// 	}

// 	// check the response
// 	log.Debug("parsed reply:", parsed_reply)
// 	switch parsed_reply.Code {
// 	case 0:
// 		// Got a reply so we can save this
// 		return nil
// 	default:
// 		// diff from ok, so it`s wrong?
// 		return fmt.Errorf("Host %s is unreachable\n", target)
// 	}
// }
