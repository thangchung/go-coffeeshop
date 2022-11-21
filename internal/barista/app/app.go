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
	"github.com/thangchung/go-coffeeshop/internal/barista/features/orders/eventhandlers"
	"github.com/thangchung/go-coffeeshop/internal/barista/features/orders/repo"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
)

type App struct {
	logger  *mylogger.Logger
	cfg     *config.Config
	network string
	address string
	handler eventhandlers.BaristaOrderedEventHandler
}

func New(log *mylogger.Logger, cfg *config.Config) *App {
	return &App{
		logger:  log,
		cfg:     cfg,
		network: "tcp",
		address: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
	}
}

func (a *App) Run() error {
	a.logger.Info("Init %s %s\n", a.cfg.Name, a.cfg.Version)

	ctx, cancel := context.WithCancel(context.Background())

	// PostgresDB
	pg, err := postgres.NewPostgresDB(a.cfg.PG.URL, postgres.MaxPoolSize(a.cfg.PG.PoolMax))
	if err != nil {
		a.logger.Fatal("app - Run - postgres.NewPostgres: %s", err.Error())

		cancel()

		return err
	}
	defer pg.Close()

	// rabbitmq
	amqpConn, err := rabbitmq.NewRabbitMQConn(a.cfg.RabbitMQ.URL, a.logger)
	if err != nil {
		cancel()

		a.logger.Fatal("app - Run - rabbitmq.NewRabbitMQConn: %s", err.Error())
	}
	defer amqpConn.Close()

	// publishers
	counterOrderPub, err := publisher.NewPublisher(
		amqpConn,
		a.logger,
		publisher.ExchangeName("counter-order-exchange"),
		publisher.BindingKey("counter-order-routing-key"),
		publisher.MessageTypeName("counter-order-updated"),
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
		a.logger,
		consumer.ExchangeName("barista-order-exchange"),
		consumer.QueueName("barista-order-queue"),
		consumer.BindingKey("barista-order-routing-key"),
		consumer.ConsumerTag("barista-order-consumer"),
	)

	if err != nil {
		a.logger.Fatal("app - Run - consumer.NewOrderConsumer: %s", err.Error())
		cancel()
	}

	go func() {
		err := consumer.StartConsumer(a.worker)
		if err != nil {
			a.logger.Error("StartConsumer: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		a.logger.Error("signal.Notify: %v", v)
	case done := <-ctx.Done():
		a.logger.Error("ctx.Done: %v", done)
	}

	a.logger.Info("Start server at " + a.address + " ...")

	return nil
}

func (c *App) worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		c.logger.Info("processDeliveries deliveryTag %v", delivery.DeliveryTag)
		c.logger.Info("received %s", delivery.Type)

		switch delivery.Type {
		case "barista-order-created":
			var payload event.BaristaOrdered
			err := json.Unmarshal(delivery.Body, &payload)

			if err != nil {
				c.logger.LogError(err)
			}

			err = c.handler.Handle(ctx, &payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					c.logger.Error("Err delivery.Reject: %v", err)
				}

				c.logger.Error("Failed to process delivery: %v", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					c.logger.Error("Failed to acknowledge delivery: %v", err)
				}
			}
		default:
			c.logger.Info("default")
		}
	}

	c.logger.Info("Deliveries channel closed")
}
