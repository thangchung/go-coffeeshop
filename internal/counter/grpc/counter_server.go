package grpc

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	counterRabbitMQ "github.com/thangchung/go-coffeeshop/internal/counter/rabbitmq"
	"github.com/thangchung/go-coffeeshop/internal/counter/usecase"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CounterServiceServerImpl struct {
	gen.UnimplementedCounterServiceServer
	logger                       *mylogger.Logger
	amqpConn                     *amqp.Connection
	cfg                          *config.Config
	productDomainService         domain.ProductDomainService
	queryOrderFulfillmentUseCase usecase.QueryOrderFulfillmentUseCase
	orderPublisher               counterRabbitMQ.OrderPublisher
}

func NewCounterServiceServerGrpc(
	grpcServer *grpc.Server,
	amqpConn *amqp.Connection,
	cfg *config.Config,
	log *mylogger.Logger,
	queryOrderFulfillmentUseCase usecase.QueryOrderFulfillmentUseCase,
	productDomainService domain.ProductDomainService,
	orderPublisher counterRabbitMQ.OrderPublisher,
) {
	svc := CounterServiceServerImpl{
		cfg:                          cfg,
		logger:                       log,
		amqpConn:                     amqpConn,
		queryOrderFulfillmentUseCase: queryOrderFulfillmentUseCase,
		productDomainService:         productDomainService,
		orderPublisher:               orderPublisher,
	}

	gen.RegisterCounterServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)
}

func (g *CounterServiceServerImpl) GetListOrderFulfillment(
	ctx context.Context,
	request *gen.GetListOrderFulfillmentRequest,
) (*gen.GetListOrderFulfillmentResponse, error) {
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

func (g *CounterServiceServerImpl) PlaceOrder(
	ctx context.Context,
	request *gen.PlaceOrderRequest,
) (*gen.PlaceOrderResponse, error) {
	g.logger.Info("POST: PlaceOrder")

	g.logger.Debug("request: %s", request)

	// add order
	order, err := domain.CreateOrderFrom(ctx, request, g.productDomainService, g.orderPublisher)
	if err != nil {
		return nil, errors.Wrap(err, "PlaceOrder - domain.CreateOrderFrom")
	}

	// todo: save to database
	// ...

	g.logger.Debug("order created: %s", *order)

	// publish order events
	// ch, err := g.amqpConn.Channel()
	// if err != nil {
	// 	panic(err)
	// }
	// defer ch.Close()

	// event := events.BaristaOrdered{
	// 	OrderID:    order.ID,
	// 	ItemLineID: uuid.New(), //todo
	// 	ItemType:   1,          //todo
	// }

	// eventBytes, err := json.Marshal(event)
	// if err != nil {
	// 	g.logger.LogError(err)
	// }

	// err = g.orderPublisher.Publish(ctx, eventBytes, "text/plain")
	// if err != nil {
	// 	g.logger.LogError(err)

	// 	return nil, errors.Wrap(err, "orderPublisher - Publish")
	// }

	// err = ch.PublishWithContext(
	// 	ctx,
	// 	OrderTopic,
	// 	"log.INFO",
	// 	false,
	// 	false,
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Type:        "barista.ordered",
	// 		Body:        eventBytes,
	// 	},
	// )

	// if err != nil {
	// 	g.logger.LogError(err)

	// 	return nil, err
	// }

	// g.logger.Info("Sending message: %s -> %s", event, g.cfg.RabbitMQ.Exchange)

	res := gen.PlaceOrderResponse{}

	return &res, nil
}
