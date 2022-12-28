package app

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/eventhandlers"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	pkgConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
)

type App struct {
	Cfg      *config.Config
	PG       postgres.DBEngine
	AMQPConn *amqp.Connection

	CounterOrderPub pkgPublisher.EventPublisher
	Consumer        pkgConsumer.EventConsumer

	handler eventhandlers.BaristaOrderedEventHandler
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	amqpConn *amqp.Connection,
	counterOrderPub pkgPublisher.EventPublisher,
	consumer pkgConsumer.EventConsumer,
	handler eventhandlers.BaristaOrderedEventHandler,
) *App {
	return &App{
		Cfg:      cfg,
		PG:       pg,
		AMQPConn: amqpConn,

		CounterOrderPub: counterOrderPub,
		Consumer:        consumer,

		handler: handler,
	}
}

func (c *App) Worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "barista-order-created":
			var payload event.BaristaOrdered

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal", err)
			}

			err = c.handler.Handle(ctx, payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					slog.Error("failed to delivery.Reject", err)
				}

				slog.Error("failed to process delivery", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					slog.Error("failed to acknowledge delivery", err)
				}
			}
		default:
			slog.Info("default")
		}
	}

	slog.Info("Deliveries channel closed")
}
