package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/features/orders/eventhandlers"
	baristaRabbitMQ "github.com/thangchung/go-coffeeshop/internal/barista/rabbitmq"
	mylog "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
)

// const (
// 	RetryTimes     = 5
// 	BackOffSeconds = 2
// )

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		glog.Fatal(err)
	}

	logger := mylog.New(cfg.Level)

	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg.RabbitMQ.URL, logger)
	if err != nil {
		logger.Fatal("app - Run - rabbitmq.NewRabbitMQConn: %s", err.Error())
	}
	defer amqpConn.Close()

	handler := eventhandlers.NewBaristaOrderedEventHandler()
	consumer, err := baristaRabbitMQ.NewOrderConsumer(amqpConn, handler, logger)

	if err != nil {
		logger.Fatal("app - Run - baristaRabbitMQ.NewOrderConsumer: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := consumer.StartConsumer(cfg.RabbitMQ.WorkerPoolSize, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.Queue, cfg.RabbitMQ.RoutingKey, cfg.RabbitMQ.ConsumerTag)
		if err != nil {
			logger.Error("StartConsumer: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.Error("signal.Notify: %v", v)
	case done := <-ctx.Done():
		logger.Error("ctx.Done: %v", done)
	}

	// rabbitConn, err := connect(cfg)
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	// defer rabbitConn.Close()

	// log.Println("Listening for and consuming RabbitMQ messages...")

	// consumer, err := event.NewConsumer(rabbitConn)
	// if err != nil {
	// 	panic(err)
	// }

	// err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	// if err != nil {
	// 	log.Println(err)
	// }
}

// func connect(cfg *config.Config) (*amqp.Connection, error) {
// 	var (
// 		counts     int64
// 		backOff    = 1 * time.Second
// 		connection *amqp.Connection
// 		rabbitURL  = cfg.RabbitMQ.URL
// 	)

// 	for {
// 		c, err := amqp.Dial(rabbitURL)
// 		if err != nil {
// 			fmt.Println("RabbitMQ not yet ready...")
// 			counts++
// 		} else {
// 			connection = c
// 			fmt.Println()

// 			break
// 		}

// 		if counts > RetryTimes {
// 			fmt.Println(err)

// 			return nil, err
// 		}

// 		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), BackOffSeconds)))
// 		backOff = time.Duration(math.Pow(float64(counts), BackOffSeconds)) * time.Second
// 		time.Sleep(backOff)

// 		continue
// 	}

// 	return connection, nil
// }
