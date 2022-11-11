package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/app"
	mylog "github.com/thangchung/go-coffeeshop/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		glog.Fatal(err)
	}

	logger := mylog.New(cfg.Level)

	a := app.New(logger, cfg)
	if err = a.Run(); err != nil {
		glog.Fatal(err)
		os.Exit(1)
	}

	// amqpConn, err := rabbitmq.NewRabbitMQConn(cfg.RabbitMQ.URL, logger)
	// if err != nil {
	// 	logger.Fatal("app - Run - rabbitmq.NewRabbitMQConn: %s", err.Error())
	// }
	// defer amqpConn.Close()

	// handler := eventhandlers.NewBaristaOrderedEventHandler()
	// consumer, err := baristaRabbitMQ.NewOrderConsumer(amqpConn, handler, logger)

	// if err != nil {
	// 	logger.Fatal("app - Run - baristaRabbitMQ.NewOrderConsumer: %s", err.Error())
	// }

	// ctx, cancel := context.WithCancel(context.Background())

	// go func() {
	// 	err := consumer.StartConsumer(cfg.RabbitMQ.WorkerPoolSize, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.Queue, cfg.RabbitMQ.RoutingKey, cfg.RabbitMQ.ConsumerTag)
	// 	if err != nil {
	// 		logger.Error("StartConsumer: %v", err)
	// 		cancel()
	// 	}
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// select {
	// case v := <-quit:
	// 	logger.Error("signal.Notify: %v", v)
	// case done := <-ctx.Done():
	// 	logger.Error("ctx.Done: %v", done)
	// }
}
