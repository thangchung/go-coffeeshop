package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/features/orders/eventhandlers"
	baristaRabbitMQ "github.com/thangchung/go-coffeeshop/internal/barista/rabbitmq"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
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

	amqpConn, err := rabbitmq.NewRabbitMQConn(a.cfg.RabbitMQ.URL, a.logger)
	if err != nil {
		a.logger.Fatal("app - Run - rabbitmq.NewRabbitMQConn: %s", err.Error())
	}
	defer amqpConn.Close()

	handler := eventhandlers.NewBaristaOrderedEventHandler()
	consumer, err := baristaRabbitMQ.NewOrderConsumer(amqpConn, handler, a.logger)

	if err != nil {
		a.logger.Fatal("app - Run - baristaRabbitMQ.NewOrderConsumer: %s", err.Error())
	}

	go func() {
		err := consumer.StartConsumer(
			a.cfg.RabbitMQ.WorkerPoolSize,
			a.cfg.RabbitMQ.Exchange,
			a.cfg.RabbitMQ.Queue,
			a.cfg.RabbitMQ.RoutingKey,
			a.cfg.RabbitMQ.ConsumerTag,
		)
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
