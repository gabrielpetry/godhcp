package log

import (
	"fmt"
	Config "go-dhcpdump/config"
	"log"
	"os"
)

var logLevels = map[string]int{
	"FATAL": 0,
	"INFO":  1,
	"ERROR": 1,
	"DEBUG": 5,
}

var config = Config.GetInstance()

func logToFile(message string) {
	logWritter, err := os.OpenFile(config.Log.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logWritter.Close()

	log.SetOutput(logWritter)
	log.Println(message)
}

func logger(logType string, v ...interface{}) {
	if logLevels[logType] > logLevels[config.Log.Level] {
		return
	}
	message := fmt.Sprintf("%s %s", logType, v)

	log.Println(message)

	if config.Log.Path != "" {
		logToFile(message)
	}

}

func Info(v ...interface{}) {
	logger("INFO", v...)
}

func Error(v ...interface{}) {
	logger("ERROR", v...)
}

func Debug(v ...interface{}) {
	logger("DEBUG", v...)
}

func Fatal(v ...interface{}) {
	logger("FATAL", v...)
	panic(v)
}
