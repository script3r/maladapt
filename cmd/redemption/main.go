package main

import (
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})

	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "redemption")
	if err != nil {
		log.Fatal("could not setup syslog logger")
	}
	log.AddHook(hook)

}
