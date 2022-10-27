package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/golang/glog"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/cmd/barista/event"
)

const (
	RetryTimes     = 5
	BackOffSeconds = 2
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		glog.Fatal(err)
	}

	rabbitConn, err := connect(cfg)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect(cfg *config.Config) (*amqp.Connection, error) {
	var (
		counts     int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
		rabbitURL  = cfg.RabbitMQ.URL
	)

	for {
		c, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			fmt.Println()

			break
		}

		if counts > RetryTimes {
			fmt.Println(err)

			return nil, err
		}

		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), BackOffSeconds)))
		backOff = time.Duration(math.Pow(float64(counts), BackOffSeconds)) * time.Second
		time.Sleep(backOff)

		continue
	}

	return connection, nil
}
