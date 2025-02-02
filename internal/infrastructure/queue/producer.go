package queue

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"time"

	"github.com/streadway/amqp"
)

// ReminderMessage — структура, которую будем отправлять в Rabbit
type ReminderMessage struct {
	UserID int64  `json:"user_id"`
	Action string `json:"action"` // SendReminder
}

type Producer interface {
	PublishReminder(msg ReminderMessage, delay time.Duration) error
}

type rabbitProducer struct {
	channel      *amqp.Channel
	exchangeName string
	exchangeType string
	routingKey   string
	log          *zerolog.Logger
}

func NewRabbitProducer(ch *amqp.Channel, exchange, routingKey string, logger *zerolog.Logger) (Producer, error) {
	err := ch.ExchangeDeclare(
		exchange,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{"x-delayed-type": "direct"},
	)
	if err != nil {
		return nil, err
	}

	return &rabbitProducer{
		channel:      ch,
		exchangeName: exchange,
		exchangeType: "x-delayed-message",
		routingKey:   routingKey,
		log:          logger,
	}, nil
}

func (p *rabbitProducer) PublishReminder(msg ReminderMessage, delay time.Duration) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	headers := amqp.Table{
		"x-delay": int64(delay / time.Millisecond),
	}

	return p.channel.Publish(
		p.exchangeName,
		p.routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:      headers,
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
