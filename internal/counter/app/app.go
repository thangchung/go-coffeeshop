package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
)

type App struct {
	logger  *mylogger.Logger
	cfg     *config.Config
	network string
	address string
}

type CounterServiceServerImpl struct {
	gen.UnimplementedCounterServiceServer
	logger     *mylogger.Logger
	rabbitConn *amqp.Connection
}

func (g *CounterServiceServerImpl) GetListOrderFulfillment(ctx context.Context, request *gen.GetListOrderFulfillmentRequest) (*gen.GetListOrderFulfillmentResponse, error) {
	g.logger.Info("GET: GetListOrderFulfillment")

	ch, err := g.rabbitConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	event := "test"

	err = ch.Publish(
		"orders_topic",
		"log.INFO",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Sending message: %s -> %s", event, "orders_topic")

	res := gen.GetListOrderFulfillmentResponse{}

	return &res, nil
}

func (g *CounterServiceServerImpl) PlaceOrder(ctx context.Context, request *gen.PlaceOrderRequest) (*gen.PlaceOrderResponse, error) {
	g.logger.Info("POST: PlaceOrder")

	res := gen.PlaceOrderResponse{}

	return &res, nil
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
	// ...

	// Use case
	// ...

	// RabbitMQ
	conn, err := connectToRabbit()
	if err != nil {
		return err
	}
	defer conn.Close()

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
	gen.RegisterCounterServiceServer(s, &CounterServiceServerImpl{logger: a.logger, rabbitConn: conn})

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	a.logger.Info("Start server at " + a.address + " ...")

	return s.Serve(l)
}

func connectToRabbit() (*amqp.Connection, error) {
	var (
		rabbitConn *amqp.Connection
		counts     int64
		rabbitURL  = "amqp://guest:guest@172.28.177.17:5672/"
	)

	for {
		connection, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Printf("rabbitmq at %s not ready...\n", rabbitURL)
			counts++
		} else {
			fmt.Println()
			rabbitConn = connection

			break
		}

		if counts > 5 {
			fmt.Println(err)

			return nil, errors.New("cannot connect to rabbit")
		}
		fmt.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)

		continue
	}
	fmt.Println("Connected to RabbitMQ!")

	return rabbitConn, nil
}
