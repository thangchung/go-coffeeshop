package consumer

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type worker func(ctx context.Context, messages <-chan amqp.Delivery)

type EventConsumer interface {
	Configure(...Option) EventConsumer
	StartConsumer(fn worker) error
}
