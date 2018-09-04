package rabbit

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Producer struct {
	rabbit *Rabbit
}

func NewProducer(rabbit *Rabbit) *Producer {
	return &Producer{
		rabbit: rabbit,
	}
}

func (p *Producer) Publish(message *ScanMessage) error {
	logger := log.WithFields(log.Fields{"func": "Publish"})
	logger.Info()

	body, err := json.Marshal(message)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = p.publish(body)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (p *Producer) Connect() error {
	logger := log.WithFields(log.Fields{"func": "Connect"})

	c, err := amqp.Dial(p.rabbit.config.AmqpUrl)
	if err != nil {
		logger.Error(err)
		return err
	}

	//Save Connection
	p.rabbit.conn = c

	//Save Channel
	p.rabbit.ch, err = p.rabbit.conn.Channel()
	if err != nil {
		logger.Error(err)
		return err
	}

	go p.DisconnectListener(p.rabbit.ch.NotifyClose(make(chan *amqp.Error)))

	if err := p.rabbit.ch.ExchangeDeclare(
		p.rabbit.config.Exchange,
		p.rabbit.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (p *Producer) DisconnectListener(ch chan *amqp.Error) {
	logger := log.WithFields(log.Fields{"func": "DisconnectListener"})
	logger.Info("Listening for disconnects..")

	err := <-ch
	if err != nil {
		logger.Info("rabbitmq disconnected")
	}
	logger.Info("rabbitmq connection closed")

	logger.Info("Attempting Reconnect..")
	p.rabbit.ch.Close()
	p.rabbit.conn.Close()

	//Reconnect
	for e := p.Connect(); e != nil; e = p.Connect() {
		logger.Info("Reconnecting..")
		time.Sleep(3 * time.Second)
	}
}

func (p *Producer) publish(msg interface{}) error {
	logger := log.WithFields(log.Fields{"func": "Publish"})

	err := p.rabbit.ch.Publish(
		p.rabbit.config.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg.([]byte),
		})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
