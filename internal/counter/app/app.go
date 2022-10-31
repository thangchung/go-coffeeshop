package app

import (
	"context"
	"errors"
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	counterGrpc "github.com/thangchung/go-coffeeshop/internal/counter/grpc"
	"github.com/thangchung/go-coffeeshop/internal/counter/usecase"
	"github.com/thangchung/go-coffeeshop/internal/counter/usecase/repo"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"google.golang.org/grpc"
)

const (
	OrderTopic     = "orders_topic"
	RetryTimes     = 5
	BackOffSeconds = 2
)

var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

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
		address: "0.0.0.0:5002",
	}
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Init %s %s\n", a.cfg.Name, a.cfg.Version)

	// Repository
	pg, err := postgres.NewPostgres(a.cfg.PG.URL, postgres.MaxPoolSize(a.cfg.PG.PoolMax))
	if err != nil {
		a.logger.Fatal("app - Run - postgres.NewPostgres: %w", err)
	}
	defer pg.Close()

	// Use case
	queryOrderFulfillmentUseCase := usecase.NewQueryOrderFulfillmentUseCase(repo.NewQueryOrderFulfillmentRepo(pg))

	// RabbitMQ
	amqpConn, err := a.connectToRabbit()
	if err != nil {
		return err
	}
	defer amqpConn.Close()

	// gRPC Server
	l, err := net.Listen(a.network, a.address)
	if err != nil {
		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			a.logger.Error("Failed to close %s %s: %v", a.network, a.address, err)
		}
	}()

	s := grpc.NewServer()
	counterGrpc.NewCounterServiceServerGrpc(s, amqpConn, queryOrderFulfillmentUseCase, a.logger)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	a.logger.Info("Start server at " + a.address + " ...")

	return s.Serve(l)
}

func (a *App) connectToRabbit() (*amqp.Connection, error) {
	var (
		amqpConn    *amqp.Connection
		counts      int64
		rabbitMqURL = a.cfg.RabbitMQ.URL
	)

	for {
		connection, err := amqp.Dial(rabbitMqURL)
		if err != nil {
			a.logger.Error("RabbitMq at %s not ready...\n", rabbitMqURL)
			counts++
		} else {
			amqpConn = connection

			break
		}

		if counts > RetryTimes {
			a.logger.LogError(err)

			return nil, ErrCannotConnectRabbitMQ
		}

		a.logger.Info("Backing off for 2 seconds...")
		time.Sleep(BackOffSeconds * time.Second)

		continue
	}

	a.logger.Info("Connected to RabbitMQ!")

	return amqpConn, nil
}
