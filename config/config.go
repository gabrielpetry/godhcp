package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Status struct {
		OnState  string
		OffState string
	}
	Log struct {
		Level string
		Path  string
	}
	Dhcpdump struct {
		Interface string
		Debug     bool
	}
	Database struct {
		Path string
	}
	Mqtt struct {
		Host   string
		Port   int
		User   string
		Passwd string
	}
}

var singleInstance *Config

var lock = &sync.Mutex{}
var systemDir = "/etc/go-dhcpdump"

func init() {
	os.MkdirAll(systemDir, os.ModeDir)
}

func setDefaults() {
	viper.SetDefault("Status", map[string]string{
		"OnState":  "ON",
		"OffState": "OFF",
	})

	viper.SetDefault("Log", map[string]string{
		"Level": "Info",
		"Path":  systemDir + "/logs/go-dhcpdump.log",
	})

	viper.SetDefault("Dhcpdump", map[string]string{
		"Interface": "en0",
		"Debug":     "false",
	})

	viper.SetDefault("Database", map[string]string{
		"Database": systemDir + "/database",
	})
}

func parseConfig() *Config {
	// viper config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/go-dhcpdump/")
	viper.AddConfigPath("$HOME/.config/go-dhcpdump")
	viper.AddConfigPath(".")

	setDefaults()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	var configuration Config
	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return &configuration
}

func GetInstance() *Config {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singleInstance == nil {
			singleInstance = parseConfig()
		}
	}
	return singleInstance
}
