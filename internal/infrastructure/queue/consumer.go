package queue

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"log"

	"github.com/streadway/amqp"
)

type Consumer interface {
	StartConsume(queueName string, handler func(ReminderMessage) error) error
}

type rabbitConsumer struct {
	channel *amqp.Channel
	log     *zerolog.Logger
}

func NewRabbitConsumer(ch *amqp.Channel, logger *zerolog.Logger) Consumer {
	return &rabbitConsumer{
		channel: ch,
		log:     logger,
	}
}

func (c *rabbitConsumer) StartConsume(queueName string, handler func(ReminderMessage) error) error {
	msgs, err := c.channel.Consume(
		queueName,
		"",
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var rm ReminderMessage
			if err := json.Unmarshal(d.Body, &rm); err != nil {
				log.Printf("Error unmarshaling: %v\n", err)
				continue
			}

			if err := handler(rm); err != nil {
				log.Printf("Handler error: %v\n", err)
			}
		}
	}()
	return nil
}
