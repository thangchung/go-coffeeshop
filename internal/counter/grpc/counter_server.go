package grpc

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/internal/counter/features"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CounterServiceServerImpl struct {
	gen.UnimplementedCounterServiceServer
	logger                  *mylogger.Logger
	amqpConn                *amqp.Connection
	cfg                     *config.Config
	productDomainSvc        domain.ProductDomainService
	queryOrderFulfillmentUC features.QueryOrderFulfillmentUseCase
	baristaOrderPub         publisher.Publisher
	kitchenOrderPub         publisher.Publisher
}

func NewCounterServiceServerGrpc(
	grpcServer *grpc.Server,
	amqpConn *amqp.Connection,
	cfg *config.Config,
	log *mylogger.Logger,
	queryOrderFulfillmentUC features.QueryOrderFulfillmentUseCase,
	productDomainSvc domain.ProductDomainService,
	baristaOrderPub publisher.Publisher,
	kitchenOrderPub publisher.Publisher,
) {
	svc := CounterServiceServerImpl{
		cfg:                     cfg,
		logger:                  log,
		amqpConn:                amqpConn,
		queryOrderFulfillmentUC: queryOrderFulfillmentUC,
		productDomainSvc:        productDomainSvc,
		baristaOrderPub:         baristaOrderPub,
		kitchenOrderPub:         kitchenOrderPub,
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

	entities, err := g.queryOrderFulfillmentUC.GetListOrderFulfillment()
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
	order, err := domain.CreateOrderFrom(ctx, request, g.productDomainSvc, g.baristaOrderPub, g.kitchenOrderPub)
	if err != nil {
		return nil, errors.Wrap(err, "PlaceOrder - domain.CreateOrderFrom")
	}

	// todo: save to database
	// ...

	g.logger.Debug("order created: %s", *order)

	res := gen.PlaceOrderResponse{}

	return &res, nil
}
