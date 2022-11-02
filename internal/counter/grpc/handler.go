package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	events "github.com/thangchung/go-coffeeshop/pkg/event"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	OrderTopic = "orders_topic"
)

func NewCounterServiceServerGrpc(
	grpcServer *grpc.Server,
	amqpConn *amqp.Connection,
	queryOrderFulfillmentUseCase domain.QueryOrderFulfillmentUseCase,
	productServiceClient domain.ProductServiceClient,
	log *mylogger.Logger) {
	svc := CounterServiceServerImpl{
		logger:                       log,
		amqpConn:                     amqpConn,
		queryOrderFulfillmentUseCase: queryOrderFulfillmentUseCase,
		productServiceClient:         productServiceClient,
	}

	gen.RegisterCounterServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)
}

type CounterServiceServerImpl struct {
	gen.UnimplementedCounterServiceServer
	logger                       *mylogger.Logger
	amqpConn                     *amqp.Connection
	productServiceClient         domain.ProductServiceClient
	queryOrderFulfillmentUseCase domain.QueryOrderFulfillmentUseCase
}

func (g *CounterServiceServerImpl) GetListOrderFulfillment(ctx context.Context, request *gen.GetListOrderFulfillmentRequest) (*gen.GetListOrderFulfillmentResponse, error) {
	g.logger.Info("GET: GetListOrderFulfillment")

	res := gen.GetListOrderFulfillmentResponse{}

	entities, err := g.queryOrderFulfillmentUseCase.GetListOrderFulfillment()
	if err != nil {
		return nil, fmt.Errorf("CounterServiceServerImpl - GetListOrderFulfillment - g.queryOrderFulfillmentUseCase.GetListOrderFulfillment: %w", err)
	}

	for _, entity := range entities {
		res.Orders = append(res.Orders, &gen.OrderDto{
			Id: entity.Id,
		})
	}

	return &res, nil
}

func (g *CounterServiceServerImpl) PlaceOrder(ctx context.Context, request *gen.PlaceOrderRequest) (*gen.PlaceOrderResponse, error) {
	g.logger.Info("POST: PlaceOrder")

	g.logger.Debug("request: %s", request)

	// add order
	order, err := domain.CreateOrderFrom(request, g.productServiceClient)
	if err != nil {
		return nil, err
	}

	g.logger.Debug("order created: %s", *order)

	// publish order events
	ch, err := g.amqpConn.Channel()
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

	g.logger.Info("Sending message: %s -> %s", event, OrderTopic)

	res := gen.PlaceOrderResponse{}

	return &res, nil
}
