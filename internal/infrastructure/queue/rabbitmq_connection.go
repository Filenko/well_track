package queue

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"well_track/internal/config"
)

// TODO:  Consider using the github.com/rabbitmq/amqp091-go package instead.

func NewRabbitMQConnection(cfg *config.RabbitMQ, log *zerolog.Logger) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.User, cfg.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	log.Info().Str("url", url).Msg("Connected to RabbitMQ")
	return conn, nil
}
