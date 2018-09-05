package rabbit

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consumer struct {
	rabbit    *Rabbit
	queueName string
}

func NewConsumer(rabbit *Rabbit) *Consumer {
	return &Consumer{
		rabbit: rabbit,
	}
}
func (c *Consumer) Connect(queueName string) error {
	logger := log.WithFields(log.Fields{"func": "ConsumerConnect"})

	con, err := amqp.Dial(c.rabbit.config.AmqpUrl)
	if err != nil {
		logger.Error(err)
		return err
	}

	//Save Connection
	c.rabbit.conn = con
	c.queueName = queueName

	//Save Channel
	c.rabbit.ch, err = c.rabbit.conn.Channel()
	if err != nil {
		logger.Error(err)
		return err
	}

	go c.DisconnectListener(c.rabbit.ch.NotifyClose(make(chan *amqp.Error)), queueName)

	if err := c.rabbit.ch.ExchangeDeclare(
		c.rabbit.config.Exchange,
		c.rabbit.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.Error(err)
		return err
	}

	q, err := c.rabbit.ch.QueueDeclare(
		queueName,
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

	err = c.rabbit.ch.QueueBind(
		q.Name,
		"",
		c.rabbit.config.Exchange,
		false,
		nil,
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Consumer) DisconnectListener(ch chan *amqp.Error, queueName string) {
	logger := log.WithFields(log.Fields{"func": "DisconnectListener"})
	logger.Info("Listening for disconnects..")

	err := <-ch
	if err != nil {
		logger.Info("rabbitmq disconnected")
	}
	logger.Info("rabbitmq connection closed")

	logger.Info("Attempting Reconnect..")
	c.rabbit.ch.Close()
	c.rabbit.conn.Close()

	//Reconnect
	c.Connect(queueName)
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	logger := log.WithFields(log.Fields{"func": "Consume"})

	deliveryChan, err := c.rabbit.ch.Consume(
		c.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return deliveryChan, nil
}
