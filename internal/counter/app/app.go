package app

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/entity"
	events "github.com/thangchung/go-coffeeshop/pkg/event"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
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

type CounterServiceServerImpl struct {
	gen.UnimplementedCounterServiceServer
	logger     *mylogger.Logger
	rabbitConn *amqp.Connection
}

type Payload struct {
	Name string `json:"name"`
}

func (g *CounterServiceServerImpl) GetListOrderFulfillment(ctx context.Context, request *gen.GetListOrderFulfillmentRequest) (*gen.GetListOrderFulfillmentResponse, error) {
	g.logger.Info("GET: GetListOrderFulfillment")

	ch, err := g.rabbitConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	event := Payload{
		Name: "drink_made",
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		g.logger.LogError(err)
	}

	err = ch.PublishWithContext(
		ctx,
		OrderTopic,
		"log.INFO",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        eventBytes,
		},
	)

	if err != nil {
		g.logger.LogError(err)

		return nil, err
	}

	g.logger.Info("Sending message: %s -> %s", event, "orders_topic")

	res := gen.GetListOrderFulfillmentResponse{}

	return &res, nil
}

func (g *CounterServiceServerImpl) PlaceOrder(ctx context.Context, request *gen.PlaceOrderRequest) (*gen.PlaceOrderResponse, error) {
	g.logger.Info("POST: PlaceOrder")

	g.logger.Debug("request: %s", request)

	// add order
	order, err := entity.CreateOrderFrom(request)
	if err != nil {
		return nil, err
	}

	g.logger.Debug("order created: %s", *order)

	// publish order events
	ch, err := g.rabbitConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	event := events.BaristaOrdered{
		OrderID:    order.ID,
		ItemLineID: uuid.New(), //todo
		ItemType:   1,          //todo
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		g.logger.LogError(err)
	}

	err = ch.PublishWithContext(
		ctx,
		OrderTopic,
		"log.INFO",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Type:        "barista.ordered",
			Body:        eventBytes,
		},
	)

	if err != nil {
		g.logger.LogError(err)

		return nil, err
	}

	g.logger.Info("Sending message: %s -> %s", event, "orders_topic")

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
	conn, err := a.connectToRabbit()
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

func (a *App) connectToRabbit() (*amqp.Connection, error) {
	var (
		rabbitConn  *amqp.Connection
		counts      int64
		rabbitMqURL = a.cfg.RabbitMQ.URL
	)

	for {
		connection, err := amqp.Dial(rabbitMqURL)
		if err != nil {
			a.logger.Error("RabbitMq at %s not ready...\n", rabbitMqURL)
			counts++
		} else {
			rabbitConn = connection

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

	return rabbitConn, nil
}
