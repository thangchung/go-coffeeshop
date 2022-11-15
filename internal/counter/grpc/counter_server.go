package grpc

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CounterServiceServerImpl struct {
	gen.UnimplementedCounterServiceServer
	logger           *mylogger.Logger
	amqpConn         *amqp.Connection
	cfg              *config.Config
	productDomainSvc domain.ProductDomainService
	orderRepo        domain.OrderRepo
	baristaOrderPub  publisher.Publisher
	kitchenOrderPub  publisher.Publisher
}

func NewCounterServiceServerGrpc(
	grpcServer *grpc.Server,
	amqpConn *amqp.Connection,
	cfg *config.Config,
	log *mylogger.Logger,
	orderRepo domain.OrderRepo,
	productDomainSvc domain.ProductDomainService,
	baristaOrderPub publisher.Publisher,
	kitchenOrderPub publisher.Publisher,
) {
	svc := CounterServiceServerImpl{
		cfg:              cfg,
		logger:           log,
		amqpConn:         amqpConn,
		orderRepo:        orderRepo,
		productDomainSvc: productDomainSvc,
		baristaOrderPub:  baristaOrderPub,
		kitchenOrderPub:  kitchenOrderPub,
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

	entities, err := g.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("CounterServiceServerImpl-GetListOrderFulfillment-g.orderRepo.GetAll: %w", err)
	}

	for _, entity := range entities {
		res.Orders = append(res.Orders, &gen.OrderDto{
			Id:              entity.ID.String(),
			OrderSource:     entity.OrderSource,
			OrderStatus:     entity.OrderStatus,
			Localtion:       entity.Location,
			LoyaltyMemberId: entity.LoyaltyMemberID.String(),
			LineItems: lo.Map(entity.LineItems, func(item *domain.LineItem, _ int) *gen.LineItemDto {
				return &gen.LineItemDto{
					ItemType:       item.ItemType,
					Name:           item.Name,
					Price:          float64(item.Price),
					ItemStatus:     item.ItemStatus,
					IsBaristaOrder: item.IsBaristaOrder,
				}
			}),
		})
	}

	return &res, nil
}

func (g *CounterServiceServerImpl) PlaceOrder(
	ctx context.Context,
	request *gen.PlaceOrderRequest,
) (*gen.PlaceOrderResponse, error) {
	g.logger.Info("POST: PlaceOrder")

	// add order
	order, err := domain.CreateOrderFrom(ctx, request, g.productDomainSvc, g.baristaOrderPub, g.kitchenOrderPub)
	if err != nil {
		return nil, errors.Wrap(err, "CounterServiceServerImpl-PlaceOrder-domain.CreateOrderFrom")
	}

	// save to database
	orderModel := &gen.OrderDto{
		Id:              order.ID.String(),
		Localtion:       order.Location,
		LoyaltyMemberId: order.LoyaltyMemberID.String(),
		OrderSource:     order.OrderSource,
		OrderStatus:     order.OrderStatus,
	}

	for _, item := range order.LineItems {
		orderModel.LineItems = append(orderModel.LineItems, &gen.LineItemDto{
			ItemType:       item.ItemType,
			Name:           item.Name,
			Price:          float64(item.Price),
			ItemStatus:     item.ItemStatus,
			IsBaristaOrder: item.IsBaristaOrder,
		})
	}

	err = g.orderRepo.Create(ctx, orderModel)
	if err != nil {
		return nil, errors.Wrap(err, "PlaceOrder-g.orderCommand.Create")
	}

	g.logger.Debug("order created: %v", *order)

	res := gen.PlaceOrderResponse{}

	return &res, nil
}
