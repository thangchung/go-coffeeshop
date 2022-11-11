package app

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	counterGrpc "github.com/thangchung/go-coffeeshop/internal/counter/grpc"
	counterRabbitMQ "github.com/thangchung/go-coffeeshop/internal/counter/rabbitmq"
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

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Init %s %s\n", a.cfg.Name, a.cfg.Version)

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

	orderPublisher, err := counterRabbitMQ.NewOrderPublisher(amqpConn, a.cfg, a.logger)
	if err != nil {
		return errors.Wrap(err, "counterRabbitMQ-NewOrderPublisher")
	}

	defer orderPublisher.CloseChan()

	a.logger.Info("Order Publisher initialized")

	var productDomainService domain.ProductDomainService = counterGrpc.NewProductServiceClient(ctx, conn)

	// Use case
	queryOrderFulfillmentUseCase := usecase.NewQueryOrderFulfillmentUseCase(ctx, repo.NewQueryOrderFulfillmentRepo(ctx, pg))

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
		queryOrderFulfillmentUseCase,
		productDomainService,
		*orderPublisher,
	)

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	a.logger.Info("Start server at " + a.address + " ...")

	return server.Serve(l)
}
