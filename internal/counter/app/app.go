package app

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/internal/counter/features/barista/eventhandlers"
	counterGrpc "github.com/thangchung/go-coffeeshop/internal/counter/grpc"
	"github.com/thangchung/go-coffeeshop/internal/counter/rabbitmq/consumer"
	counterPublisher "github.com/thangchung/go-coffeeshop/internal/counter/rabbitmq/publisher"
	"github.com/thangchung/go-coffeeshop/internal/counter/usecase"
	"github.com/thangchung/go-coffeeshop/internal/counter/usecase/repo"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	logger  *mylogger.Logger
	cfg     *config.Config
	network string
	address string
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

		return err
	}
	defer pg.Close()

	// RabbitMQ
	amqpConn, err := rabbitmq.NewRabbitMQConn(a.cfg.RabbitMQ.URL, a.logger)
	if err != nil {
		a.logger.Fatal("app - Run - rabbitmq.NewRabbitMQConn: %s", err.Error())

		return err
	}
	defer amqpConn.Close()

	// gRPC Client
	conn, err := grpc.Dial(a.cfg.ProductClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	baristaOrderPub, err := counterPublisher.NewPublisher(
		amqpConn,
		a.logger,
		counterPublisher.ExchangeName("barista-order-exchange"),
		counterPublisher.BindingKey("barista-order-routing-key"),
		counterPublisher.MessageTypeName("barista-order-created"),
	)
	defer baristaOrderPub.CloseChan()

	if err != nil {
		return errors.Wrap(err, "counterRabbitMQ-Barista-NewOrderPublisher")
	}

	kitchenOrderPub, err := counterPublisher.NewPublisher(
		amqpConn,
		a.logger,
		counterPublisher.ExchangeName("kitchen-order-exchange"),
		counterPublisher.BindingKey("kitchen-order-routing-key"),
		counterPublisher.MessageTypeName("kitchen-order-created"),
	)
	defer kitchenOrderPub.CloseChan()

	if err != nil {
		return errors.Wrap(err, "counterRabbitMQ-Kitchen-NewOrderPublisher")
	}

	a.logger.Info("Order Publisher initialized")

	var productDomainSvc domain.ProductDomainService = counterGrpc.NewProductServiceClient(ctx, conn)

	// Use case
	queryOrderFulfillmentUC := usecase.NewQueryOrderFulfillmentUseCase(ctx, repo.NewQueryOrderFulfillmentRepo(ctx, pg))

	// event handlers.
	handler := eventhandlers.NewBaristaOrderUpdatedEventHandler()

	// consumers
	consumer, err := consumer.NewConsumer(
		amqpConn,
		handler,
		a.logger,
		consumer.ExchangeName("counter-order-exchange"),
		consumer.QueueName("counter-order-queue"),
		consumer.BindingKey("counter-order-routing-key"),
		consumer.ConsumerTag("counter-order-consumer"),
		consumer.MessageTypeName("counter-order-updated"),
	)

	if err != nil {
		a.logger.Fatal("app - Run - consumer.NewOrderConsumer: %s", err.Error())
	}

	go func() {
		err = consumer.StartConsumer()
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
		queryOrderFulfillmentUC,
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
