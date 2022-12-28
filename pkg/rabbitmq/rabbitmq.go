package rabbitmq

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

type RabbitMQConnStr string

var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

func NewRabbitMQConn(rabbitMqURL RabbitMQConnStr) (*amqp.Connection, error) {
	var (
		amqpConn *amqp.Connection
		counts   int64
	)

	for {
		connection, err := amqp.Dial(string(rabbitMqURL))
		if err != nil {
			slog.Error("failed to connect to RabbitMq...", err, rabbitMqURL)
			counts++
		} else {
			amqpConn = connection

			break
		}

		if counts > _retryTimes {
			slog.Error("failed to retry", err)

			return nil, ErrCannotConnectRabbitMQ
		}

		slog.Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)

		continue
	}

	slog.Info("📫 connected to rabbitmq 🎉")

	return amqpConn, nil
}
