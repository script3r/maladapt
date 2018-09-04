package rabbit

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type ScanMessage struct {
	Filename string `json:"filename"`
	SHA256   string `json:"sha256"`
	MD5      string `json:"md5"`
	Path     string `json:"path"`
}

func NewScanMessage(filename string, hSha256 [sha256.Size]byte, hMd5 [md5.Size]byte, path string) *ScanMessage {
	return &ScanMessage{
		Filename: filename,
		SHA256:   base64.RawURLEncoding.EncodeToString(hSha256[:]),
		MD5:      base64.RawURLEncoding.EncodeToString(hMd5[:]),
		Path:     path,
	}
}

type Rabbit struct {
	config Configuration
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func (r *Rabbit) ConsumerConnect(queueName string) error {
	logger := log.WithFields(log.Fields{"func": "ConsumerConnect"})

	c, err := amqp.Dial(r.config.AmqpUrl)
	if err != nil {
		logger.Error(err)
		return err
	}

	//Save Connection
	r.conn = c

	//Save Channel
	r.ch, err = r.conn.Channel()
	if err != nil {
		logger.Error(err)
		return err
	}

	go r.ConsumerDisconnectListener(r.ch.NotifyClose(make(chan *amqp.Error)), queueName)

	if err := r.ch.ExchangeDeclare(
		r.config.Exchange,
		r.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.Error(err)
		return err
	}

	q, err := r.ch.QueueDeclare(
		"queueName",
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = r.ch.QueueBind(
		q.Name,
		"",
		r.config.Exchange,
		false,
		nil,
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r *Rabbit) ConsumerDisconnectListener(ch chan *amqp.Error, queueName string) {
	logger := log.WithFields(log.Fields{"func": "DisconnectListener"})
	logger.Info("Listening for disconnects..")

	err := <-ch
	if err != nil {
		logger.Info("rabbitmq disconnected")
	}
	logger.Info("rabbitmq connection closed")

	logger.Info("Attempting Reconnect..")
	r.ch.Close()
	r.conn.Close()

	//Reconnect
	r.ConsumerConnect(queueName)
}

func NewRabbit(config Configuration) *Rabbit {
	return &Rabbit{
		config: config,
	}
}
