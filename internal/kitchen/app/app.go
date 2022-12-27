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
	"github.com/thangchung/go-coffeeshop/cmd/kitchen/config"
	"github.com/thangchung/go-coffeeshop/internal/kitchen/eventhandlers"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	pkgRabbitMQ "github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	pkgConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
)

type App struct {
	cfg     *config.Config
	network string
	address string
	handler eventhandlers.KitchenOrderedEventHandler
}

func New(cfg *config.Config) *App {
	return &App{
		cfg:     cfg,
		network: "tcp",
		address: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
	}
}

func (a *App) Run() error {
	slog.Info("init app", "name", a.cfg.Name, "version", a.cfg.Version)

	ctx, cancel := context.WithCancel(context.Background())

	// postgresdb.
	pg, err := postgres.NewPostgresDB(postgres.DBConnString(a.cfg.PG.DsnURL))
	if err != nil {
		cancel()

		slog.Error("failed to create a new Postgres", err)

		return err
	}
	defer pg.Close()

	// rabbitmq.
	amqpConn, err := pkgRabbitMQ.NewRabbitMQConn(pkgRabbitMQ.RabbitMQConnStr(a.cfg.RabbitMQ.URL))
	if err != nil {
		cancel()

		slog.Error("failed to create a new RabbitMQConn", err)
	}
	defer amqpConn.Close()

	// publishers
	counterOrderPub, err := pkgPublisher.NewPublisher(
		amqpConn,
	)
	if err != nil {
		cancel()

		return errors.Wrap(err, "publisher-Counter-NewOrderPublisher")
	}
	defer counterOrderPub.CloseChan()

	counterOrderPub.Configure(
		pkgPublisher.ExchangeName("counter-order-exchange"),
		pkgPublisher.BindingKey("counter-order-routing-key"),
		pkgPublisher.MessageTypeName("kitchen-order-updated"),
	)

	// event handlers.
	a.handler = eventhandlers.NewKitchenOrderedEventHandler(pg, counterOrderPub)

	// consumers.
	consumer, err := pkgConsumer.NewConsumer(
		amqpConn,
	)
	if err != nil {
		slog.Error("failed to create a new OrderConsumer", err)
		cancel()
	}

	consumer.Configure(
		pkgConsumer.ExchangeName("kitchen-order-exchange"),
		pkgConsumer.QueueName("kitchen-order-queue"),
		pkgConsumer.BindingKey("kitchen-order-routing-key"),
		pkgConsumer.ConsumerTag("kitchen-order-consumer"),
	)

	go func() {
		err := consumer.StartConsumer(a.worker)
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

	slog.Info("start server...", "address", a.address)

	return nil
}

func (c *App) worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "kitchen-order-created":
			var payload event.KitchenOrdered
			err := json.Unmarshal(delivery.Body, &payload)

			if err != nil {
				slog.Error("failed to Unmarshal message", err)
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

	slog.Info("deliveries channel closed")
}
