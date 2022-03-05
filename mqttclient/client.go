package mqttclient

import (
	"fmt"
	Config "go-dhcpdump/config"
	"go-dhcpdump/log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var config = Config.GetInstance()

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info("Connected to MQTT broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Fatal(err)
}

var mqttClient *mqtt.Client
var lock = &sync.Mutex{}

func GetClient() mqtt.Client {
	if mqttClient != nil {
		return *mqttClient
	}

	lock.Lock()
	defer lock.Unlock()

	broker := config.Mqtt.Host
	port := config.Mqtt.Port
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("mqttclient")
	opts.SetUsername(config.Mqtt.User)
	opts.SetPassword(config.Mqtt.Passwd)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient = &client
	return client
}
