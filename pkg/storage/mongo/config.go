package mongo

import (
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	Hosts        []string
	Database     string
	Username     string
	Password     string
	VerifyTLS    bool
	WriteConcern int
	Timeout      int64
}

func NewConfiguration(hosts []string, database string, username string, password string, verifyTLS bool, writeConcern int, timeout int64) Configuration {
	log.WithFields(log.Fields{"func": "NewConfiguraton", "hosts": hosts, "database": database}).Info()
	return Configuration{
		Hosts:        hosts,
		Database:     database,
		Username:     username,
		Password:     password,
		VerifyTLS:    verifyTLS,
		WriteConcern: writeConcern,
		Timeout:      timeout,
	}
}
