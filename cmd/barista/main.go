package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/barista/event"
)

const (
	RetryTimes = 5
	PowOf      = 2
)

func main() {
	rabbitConn, err := connect()
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

func connect() (*amqp.Connection, error) {
	var (
		counts     int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
		rabbitURL  = "amqp://guest:guest@172.28.177.17:5672/"
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

		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), PowOf)))
		backOff = time.Duration(math.Pow(float64(counts), PowOf)) * time.Second
		time.Sleep(backOff)

		continue
	}

	return connection, nil
}
