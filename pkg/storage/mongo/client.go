package mongo

import (
	"crypto/tls"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type MongoClient struct {
	config Configuration
}

func (m *MongoClient) Session() (*mgo.Session, error) {
	logger := log.WithFields(log.Fields{"func": "Session"})

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    m.config.Hosts,
		Database: m.config.Database,
		Username: m.config.Username,
		Password: m.config.Password,
		DialServer: func(a *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", a.String(), &tls.Config{
				InsecureSkipVerify: !m.config.VerifyTLS,
			})
		},
		Timeout: time.Second * time.Duration(m.config.Timeout),
	})

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := session.Ping(); err != nil {
		logger.Error(err)

		session.Close()

		return nil, err
	}

	// explicitly set write concern
	session.EnsureSafe(&mgo.Safe{W: m.config.WriteConcern})

	return session, nil
}

func (m *MongoClient) Database() (*mgo.Session, *mgo.Database, error) {
	logger := log.WithFields(log.Fields{"func": "Database"})
	s, err := m.Session()
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return s, s.DB(m.config.Database), nil
}

func NewMongoClient(config Configuration) *MongoClient {
	return &MongoClient{config}
}
