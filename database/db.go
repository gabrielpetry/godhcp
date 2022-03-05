package database

import (
	Config "go-dhcpdump/config"
	log "go-dhcpdump/log"
	"sync"

	scribble "github.com/nanobox-io/golang-scribble"
)

var config = Config.GetInstance()
var lock = &sync.Mutex{}

var db *scribble.Driver

func GetInstance() *scribble.Driver {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()

		driver, err := scribble.New(config.Database.Path, nil)
		if err != nil {
			log.Fatal(err)
		}
		db = driver
	}

	return db
}
