package ping

import (
	"container/list"
	"go-dhcpdump/dhcpMessage"
	"go-dhcpdump/log"
	"time"
)

type Job struct {
	dhcp         *dhcpMessage.DhcpdumpMessage
	errorCounter int
	error        error
}

const errorLimit = 3

// var jobs = make(chan Job)
var jobs = list.New()

func AddJob(dhcp dhcpMessage.DhcpdumpMessage) {
	job := Job{dhcp: &dhcp, errorCounter: 0, error: nil}
	jobs.PushBack(job)
	log.Debug("Adding new job", job)
}

func retryJob(job Job) {
	jobs.PushBack(job)
}

func StartWorker() {
	log.Info("Starting worker")
	for {
		time.Sleep(time.Millisecond * 500)
		if jobs.Len() < 1 {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		// get the element and remove
		e := jobs.Front()
		jobs.Remove(e)
		// cast to Job struct
		job := e.Value.(Job)
		log.Debug("job", job)

		if job.errorCounter >= errorLimit {
			log.Debug("target stopped responding", job)
			job.dhcp.DisconnectDevice()
			continue
		}

		err := Ping(job.dhcp.ClientIpAddress)

		if err != nil {
			log.Error("Error Ping Device: ", job.dhcp, err)
			job.errorCounter += 1
		} else {
			log.Debug("Device Responded OK", job.dhcp)
			job.errorCounter = 0
		}

		retryJob(job)
	}
}
