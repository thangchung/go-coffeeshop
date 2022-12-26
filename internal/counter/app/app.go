package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/events"
	"github.com/thangchung/go-coffeeshop/internal/counter/events/handlers"
	counterGrpc "github.com/thangchung/go-coffeeshop/internal/counter/infras/grpc"
	"github.com/thangchung/go-coffeeshop/internal/counter/infras/repo"
	"github.com/thangchung/go-coffeeshop/internal/counter/usecases/orders"
	sharedevents "github.com/thangchung/go-coffeeshop/internal/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	rabConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg            *config.Config
	network        string
	address        string
	baristaHandler events.BaristaOrderUpdatedEventHandler
	kitchenHandler events.KitchenOrderUpdatedEventHandler
}

func New(cfg *config.Config) *App {
	return &App{
		cfg:     cfg,
		network: "tcp",
		address: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
	}
}

func (a *App) Run() error {
	slog.Info("Init app", "name", a.cfg.Name, "version", a.cfg.Version)

	ctx, cancel := context.WithCancel(context.Background())

	// postgresdb.
	pg, err := postgres.NewPostgresDB(a.cfg.PG.DsnURL)
	if err != nil {
		cancel()

		slog.Error("failed to create a new Postgres", err)

		return err
	}
	defer pg.Close()

	// rabbitmq.
	amqpConn, err := rabbitmq.NewRabbitMQConn(a.cfg.RabbitMQ.URL)
	if err != nil {
		slog.Error("failed to create a new RabbitMQConn", err)

		cancel()

		return err
	}
	defer amqpConn.Close()

	// gRPC Client.
	conn, err := grpc.Dial(a.cfg.ProductClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()

		return err
	}
	defer conn.Close()

	baristaOrderPub, err := publisher.NewPublisher(
		amqpConn,
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
		publisher.ExchangeName("kitchen-order-exchange"),
		publisher.BindingKey("kitchen-order-routing-key"),
		publisher.MessageTypeName("kitchen-order-created"),
	)
	defer kitchenOrderPub.CloseChan()

	if err != nil {
		cancel()

		return errors.Wrap(err, "counterRabbitMQ-Kitchen-NewOrderPublisher")
	}

	slog.Info("Order Publisher initialized")

	// repository
	orderRepo := repo.NewOrderRepo(pg)

	// domain service
	productDomainSvc := counterGrpc.NewGRPCProductClient(conn)

	// usecases.
	uc := orders.NewUseCase(
		orderRepo,
		productDomainSvc,
		baristaOrderPub,
		kitchenOrderPub,
	)

	// event handlers.
	a.baristaHandler = handlers.NewBaristaOrderUpdatedEventHandler(orderRepo)
	a.kitchenHandler = handlers.NewKitchenOrderUpdatedEventHandler(orderRepo)

	// consumers
	consumer, err := rabConsumer.NewConsumer(
		amqpConn,
		rabConsumer.ExchangeName("counter-order-exchange"),
		rabConsumer.QueueName("counter-order-queue"),
		rabConsumer.BindingKey("counter-order-routing-key"),
		rabConsumer.ConsumerTag("counter-order-consumer"),
	)

	if err != nil {
		slog.Error("failed to create a new consumer", err)
	}

	go func() {
		err = consumer.StartConsumer(a.worker)
		if err != nil {
			slog.Error("failed to start consumer: %v", err)
			cancel()
		}
	}()

	// gRPC Server
	l, err := net.Listen(a.network, a.address)
	if err != nil {
		slog.Error("failed to listen to address", err, "network", a.network, "address", a.address)

		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			slog.Error("failed to close", err, "network", a.network, "address", a.address)
		}
	}()

	server := grpc.NewServer()
	counterGrpc.NewGRPCCounterServer(
		server,
		a.cfg,
		uc,
	)

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	slog.Info("start server...", "address", a.address)

	return server.Serve(l)
}

func (c *App) worker(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "barista-order-updated":
			var payload sharedevents.BaristaOrderUpdated

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}

			err = c.baristaHandler.Handle(ctx, &payload)

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
		case "kitchen-order-updated":
			var payload sharedevents.KitchenOrderUpdated

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}

			err = c.kitchenHandler.Handle(ctx, &payload)

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
