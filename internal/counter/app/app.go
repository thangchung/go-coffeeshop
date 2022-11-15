package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/internal/counter/features/orders/eventhandlers"
	"github.com/thangchung/go-coffeeshop/internal/counter/features/orders/repo"
	counterGrpc "github.com/thangchung/go-coffeeshop/internal/counter/grpc"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	rabConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	logger  *mylogger.Logger
	cfg     *config.Config
	network string
	address string
	handler domain.BaristaOrderUpdatedEventHandler
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

	// RabbitMQ
	amqpConn, err := rabbitmq.NewRabbitMQConn(a.cfg.RabbitMQ.URL, a.logger)
	if err != nil {
		a.logger.Fatal("app - Run - rabbitmq.NewRabbitMQConn: %s", err.Error())

		cancel()

		return err
	}
	defer amqpConn.Close()

	// gRPC Client
	conn, err := grpc.Dial(a.cfg.ProductClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()

		return err
	}
	defer conn.Close()

	baristaOrderPub, err := publisher.NewPublisher(
		amqpConn,
		a.logger,
		publisher.ExchangeName("barista-order-exchange"),
		publisher.BindingKey("barista-order-routing-key"),
		publisher.MessageTypeName("barista-order-created"),
	)
	defer baristaOrderPub.CloseChan()

	if err != nil {
		cancel()

		return errors.Wrap(err, "counterRabbitMQ-Barista-NewOrderPublisher")
	}

	kitchenOrderPub, err := publisher.NewPublisher(
		amqpConn,
		a.logger,
		publisher.ExchangeName("kitchen-order-exchange"),
		publisher.BindingKey("kitchen-order-routing-key"),
		publisher.MessageTypeName("kitchen-order-created"),
	)
	defer kitchenOrderPub.CloseChan()

	if err != nil {
		cancel()

		return errors.Wrap(err, "counterRabbitMQ-Kitchen-NewOrderPublisher")
	}

	a.logger.Info("Order Publisher initialized")

	// repository
	orderRepo := repo.NewOrderRepo(pg)

	// domain service
	productDomainSvc := counterGrpc.NewProductDomainService(conn)

	// event handlers.
	a.handler = eventhandlers.NewBaristaOrderUpdatedEventHandler(orderRepo)

	// consumers
	consumer, err := rabConsumer.NewConsumer(
		amqpConn,
		a.logger,
		rabConsumer.ExchangeName("counter-order-exchange"),
		rabConsumer.QueueName("counter-order-queue"),
		rabConsumer.BindingKey("counter-order-routing-key"),
		rabConsumer.ConsumerTag("counter-order-consumer"),
	)

	if err != nil {
		a.logger.Fatal("app - Run - consumer.NewOrderConsumer: %s", err.Error())
	}

	go func() {
		err = consumer.StartConsumer(a.worker)
		if err != nil {
			a.logger.Error("StartConsumer: %v", err)
			cancel()
		}
	}()

	// gRPC Server
	l, err := net.Listen(a.network, a.address)
	if err != nil {
		a.logger.Fatal("app - Run - net.Listener: %s", err.Error())

		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			a.logger.Error("Failed to close %s %s: %v", a.network, a.address, err)
		}
	}()

	server := grpc.NewServer()
	counterGrpc.NewCounterServiceServerGrpc(
		server,
		amqpConn,
		a.cfg,
		a.logger,
		orderRepo,
		productDomainSvc,
		*baristaOrderPub,
		*kitchenOrderPub,
	)

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	a.logger.Info("Start server at " + a.address + " ...")

	return server.Serve(l)
}

func (c *App) worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		c.logger.Info("processDeliveries deliveryTag %v", delivery.DeliveryTag)
		c.logger.Info("received %s", delivery.Type)

		switch delivery.Type {
		case "counter-order-updated":
			var payload event.BaristaOrderUpdated
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
