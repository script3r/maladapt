package rabbit

import (
	"log"
	"testing"
)

func TestConsumer(t *testing.T) {
	c := NewConsumer(
		NewRabbit(
			NewConfiguration("amqp://localhost:",
				"maladapt",
				"fanout",
				5),
		),
	)

	c.Connect("test_queue")

	msgs, err := c.Consume()
	if err != nil {
		log.Fatal(err)
	}

	d := <-msgs
	log.Printf(" [x] %s", d.Body)
}

func TestMultipleConsumers(t *testing.T) {
	c := NewConsumer(
		NewRabbit(
			NewConfiguration("amqp://localhost:",
				"maladapt",
				"fanout",
				5),
		),
	)

	c.Connect("test_queue")

	c2 := NewConsumer(
		NewRabbit(
			NewConfiguration("amqp://localhost:",
				"maladapt",
				"fanout",
				5),
		),
	)

	c2.Connect("test_queue2")
	msgs, err := c.Consume()
	if err != nil {
		log.Fatal(err)
	}
	msgs2, err := c2.Consume()
	if err != nil {
		log.Fatal(err)
	}

	d := <-msgs
	log.Printf(" [x] %s", d.Body)
	d2 := <-msgs2
	log.Printf(" [x] %s", d2.Body)
}
