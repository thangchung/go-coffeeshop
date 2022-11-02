package rabbitmq

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
)

const (
	OrderTopic     = "orders_topic"
	RetryTimes     = 5
	BackOffSeconds = 2
)

var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

func NewRabbitMQConn(rabbitMqURL string, logger *mylogger.Logger) (*amqp.Connection, error) {
	var (
		amqpConn *amqp.Connection
		counts   int64
	)

	for {
		connection, err := amqp.Dial(rabbitMqURL)
		if err != nil {
			logger.Error("RabbitMq at %s not ready...\n", rabbitMqURL)
			counts++
		} else {
			amqpConn = connection

			break
		}

		if counts > RetryTimes {
			logger.LogError(err)

			return nil, ErrCannotConnectRabbitMQ
		}

		logger.Info("Backing off for 2 seconds...")
		time.Sleep(BackOffSeconds * time.Second)

		continue
	}

	logger.Info("Connected to RabbitMQ!")

	return amqpConn, nil
}
