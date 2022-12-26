package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/eventhandlers"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	pkgConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
)

func Run(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) error {
	slog.Info("⚡ init app", "name", cfg.Name, "version", cfg.Version)

	// postgresdb.
	pg, err := postgres.NewPostgresDB(cfg.PG.DsnURL)
	if err != nil {
		cancel()

		slog.Error("failed to create a new Postgres", err)

		return err
	}
	defer pg.Close()

	// rabbitmq.
	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg.RabbitMQ.URL)
	if err != nil {
		cancel()

		slog.Error("failed to create a new RabbitMQConn", err)

		return err
	}
	defer amqpConn.Close()

	// publishers.
	counterOrderPub, err := pkgPublisher.NewPublisher(
		amqpConn,
		pkgPublisher.ExchangeName("counter-order-exchange"),
		pkgPublisher.BindingKey("counter-order-routing-key"),
		pkgPublisher.MessageTypeName("barista-order-updated"),
	)
	defer counterOrderPub.CloseChan()

	if err != nil {
		cancel()

		return errors.Wrap(err, "publisher-Counter-NewOrderPublisher")
	}

	// consumers.
	consumer, err := pkgConsumer.NewConsumer(
		amqpConn,
		pkgConsumer.ExchangeName("barista-order-exchange"),
		pkgConsumer.QueueName("barista-order-queue"),
		pkgConsumer.BindingKey("barista-order-routing-key"),
		pkgConsumer.ConsumerTag("barista-order-consumer"),
	)
	if err != nil {
		slog.Error("failed to create a new OrderConsumer", err)
		cancel()
	}

	a, err := InitApp(cfg, pg, amqpConn, counterOrderPub, consumer)
	if err != nil {
		slog.Error("failed init app", err)
		cancel()
	}

	// event handlers.
	a.handler = eventhandlers.NewBaristaOrderedEventHandler(pg, counterOrderPub)

	slog.Info("🌏 start server...", "address", fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port))

	go func() {
		err := a.consumer.StartConsumer(a.worker)
		if err != nil {
			slog.Error("failed to start Consumer", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		slog.Info("signal.Notify", v)
	case done := <-ctx.Done():
		slog.Info("ctx.Done", done)
	}

	return nil
}

type App struct {
	cfg *config.Config

	pg       *postgres.Postgres
	amqpConn *amqp.Connection

	counterOrderPub rabbitmq.EventPublisher
	consumer        *pkgConsumer.Consumer

	handler eventhandlers.BaristaOrderedEventHandler
}

func New(
	cfg *config.Config,
	pg *postgres.Postgres,
	amqpConn *amqp.Connection,
	counterOrderPub rabbitmq.EventPublisher,
	consumer *pkgConsumer.Consumer,
	handler eventhandlers.BaristaOrderedEventHandler,
) *App {
	return &App{
		cfg: cfg,

		pg:       pg,
		amqpConn: amqpConn,

		counterOrderPub: counterOrderPub,
		consumer:        consumer,

		handler: handler,
	}
}

func (c *App) worker(ctx context.Context, messages <-chan amqp091.Delivery) {
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
