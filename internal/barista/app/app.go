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
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/eventhandlers"
	"github.com/thangchung/go-coffeeshop/internal/barista/infras/repo"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
)

type App struct {
	cfg     *config.Config
	network string
	address string
	handler eventhandlers.BaristaOrderedEventHandler
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

	// PostgresDB
	// pg, err := postgres.NewPostgresDB(a.cfg.PG.URL, postgres.MaxPoolSize(a.cfg.PG.PoolMax))
	// if err != nil {
	// 	slog.Error("failed to create a new Postgres", err, err.Error())

	// 	cancel()

	// 	return err
	// }
	// defer pg.Close()

	pg, err := postgres.NewPostgreSQLDb(a.cfg.PG.URL)
	if err != nil {
		cancel()

		slog.Error("failed to create a new Postgres", err, err.Error())

		return err
	}
	defer pg.CloseDB()

	// rabbitmq
	amqpConn, err := rabbitmq.NewRabbitMQConn(a.cfg.RabbitMQ.URL)
	if err != nil {
		cancel()

		slog.Error("failed to create a new RabbitMQConn", err, err.Error())

		return err
	}
	defer amqpConn.Close()

	// publishers
	counterOrderPub, err := publisher.NewPublisher(
		amqpConn,
		publisher.ExchangeName("counter-order-exchange"),
		publisher.BindingKey("counter-order-routing-key"),
		publisher.MessageTypeName("barista-order-updated"),
	)
	defer counterOrderPub.CloseChan()

	if err != nil {
		cancel()

		return errors.Wrap(err, "publisher-Counter-NewOrderPublisher")
	}

	// repository
	orderRepo := repo.NewOrderRepo(pg)

	// event handlers.
	a.handler = eventhandlers.NewBaristaOrderedEventHandler(orderRepo, counterOrderPub)

	// consumers
	consumer, err := consumer.NewConsumer(
		amqpConn,
		consumer.ExchangeName("barista-order-exchange"),
		consumer.QueueName("barista-order-queue"),
		consumer.BindingKey("barista-order-routing-key"),
		consumer.ConsumerTag("barista-order-consumer"),
	)
	if err != nil {
		slog.Error("failed to create a new OrderConsumer", err, err.Error())
		cancel()
	}

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
		case "barista-order-created":
			var payload event.BaristaOrdered
			err := json.Unmarshal(delivery.Body, &payload)

			if err != nil {
				slog.Error("failed to Unmarshal", err, err.Error())
			}

			err = c.handler.Handle(ctx, &payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					slog.Error("failed to delivery.Reject", err, err.Error())
				}

				slog.Error("failed to process delivery", err, err.Error())
			} else {
				err = delivery.Ack(false)
				if err != nil {
					slog.Error("failed to acknowledge delivery", err, err.Error())
				}
			}
		default:
			slog.Info("default")
		}
	}

	slog.Info("Deliveries channel closed")
}
