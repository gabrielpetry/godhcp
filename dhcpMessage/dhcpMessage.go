package dhcpMessage

import (
	"encoding/json"
	"fmt"
	Config "go-dhcpdump/config"
	"go-dhcpdump/database"
	log "go-dhcpdump/log"
	"go-dhcpdump/mqttclient"
	"regexp"
	"strings"
	"sync"
	"time"
)

var config = Config.GetInstance()

var db = database.GetInstance()
var dbCollection = "dhcpclients"
var lock = &sync.Mutex{}

func (dhcp *DhcpdumpMessage) findByRegex(regex string, message *string, sep string) string {
	r, _ := regexp.Compile(regex)
	match := r.FindString(*message)

	split := strings.Split(match, sep)
	if len(split) > 1 {
		return split[1]
	}
	return split[0]
}

func (dhcp *DhcpdumpMessage) Parse(message string) {
	lock.Lock()
	defer lock.Unlock()
	log.Debug("Parsing  message:")
	message = strings.Trim(message, "\n")

	entries := strings.Split(message, "\n")
	metadata := entries[len(entries)-1]

	macRegex := "(ClientHWAddr)=[\\d\\w]{2}:[\\d\\w]{2}:[\\d\\w]{2}:[\\d\\w]{2}:[\\d\\w]{2}:[\\d\\w]{2}"
	ipRegex := "(RequestIP):\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}"
	hostnameRegex := "(Hostname):[a-zA-Z0-9-]+"

	dhcp.ClientHostName = dhcp.findByRegex(hostnameRegex, &metadata, ":")
	dhcp.ClientIpAddress = dhcp.findByRegex(ipRegex, &metadata, ":")
	dhcp.ClientMacAddress = dhcp.findByRegex(macRegex, &metadata, "=")

	if dhcp.ClientHostName == "" {
		dhcp.ClientHostName = dhcp.ClientMacAddress
	}

	dhcp.ClientStatus = config.Status.OnState
}

func (dhcp *DhcpdumpMessage) Save() error {
	if dhcp.ClientIpAddress == "" {
		return fmt.Errorf("ClientIpAddress is empty, %v", dhcp)
	}

	if err := db.Write(dbCollection, dhcp.ClientHostName, dhcp); err != nil {
		return err
	}

	dhcp.Publish()
	return nil
}

func (dhcp *DhcpdumpMessage) GetAll() []DhcpdumpMessage {
	results, err := db.ReadAll(dbCollection)

	if err != nil {
		log.Error(err)
		return []DhcpdumpMessage{}
	}

	clients := []DhcpdumpMessage{}

	for _, f := range results {
		clientFound := DhcpdumpMessage{}
		if err := json.Unmarshal([]byte(f), &clientFound); err != nil {
			fmt.Println("Error", err)
		}
		clients = append(clients, clientFound)
	}

	return clients
}

func (dhcp *DhcpdumpMessage) pub(topic string) {
	mqtt := mqttclient.GetClient()
	topic = fmt.Sprintf("home/devices/%s", topic)
	token := mqtt.Publish(topic, 0, false, dhcp.ClientStatus)
	token.Wait()
	time.Sleep(time.Second)
}

func (dhcp *DhcpdumpMessage) Publish() {
	log.Debug("Publishing message:", dhcp)
	dhcp.pub(dhcp.ClientMacAddress)
	if dhcp.ClientMacAddress != dhcp.ClientHostName {
		dhcp.pub(dhcp.ClientHostName)
	}
}

func (dhcp *DhcpdumpMessage) DisconnectDevice() {
	dhcp.ClientStatus = config.Status.OffState
	dhcp.Save()
}

func (dhcp *DhcpdumpMessage) GetOnlineDevices() []DhcpdumpMessage {
	clients := dhcp.GetAll()
	onlineDevices := []DhcpdumpMessage{}

	for _, client := range clients {
		if client.ClientStatus == config.Status.OnState {
			onlineDevices = append(onlineDevices, client)
		}
	}

	return onlineDevices
}

func (dhcp *DhcpdumpMessage) FilterDhcpPacket() error {
	if len(config.Filter.ClientHostName) == len(config.Filter.ClientMacAddress) &&
		len(config.Filter.ClientHostName) == 0 {
		return nil
	}

	for _, v := range config.Filter.ClientHostName {
		if v == dhcp.ClientHostName {
			return nil
		}
	}

	for _, v := range config.Filter.ClientMacAddress {
		if v == dhcp.ClientMacAddress {
			return nil
		}
	}

	return fmt.Errorf("Client ain't in the allowed list, %v", dhcp)
}
