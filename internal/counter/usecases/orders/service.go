package orders

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
	"golang.org/x/exp/slog"
)

type usecase struct {
	orderRepo        domain.OrderRepo
	productDomainSvc domain.ProductDomainService
	baristaEventPub  shared.EventPublisher
	kitchenEventPub  shared.EventPublisher
}

var _ UseCase = (*usecase)(nil)

func NewUseCase(
	orderRepo domain.OrderRepo,
	productDomainSvc domain.ProductDomainService,
	baristaEventPub shared.EventPublisher,
	kitchenEventPub shared.EventPublisher,
) UseCase {
	return &usecase{
		orderRepo:        orderRepo,
		productDomainSvc: productDomainSvc,
		baristaEventPub:  baristaEventPub,
		kitchenEventPub:  kitchenEventPub,
	}
}

func (uc *usecase) GetListOrderFulfillment(ctx context.Context) ([]*domain.Order, error) {
	entities, err := uc.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("counterGRPCServer-GetListOrderFulfillment-g.orderRepo.GetAll: %w", err)
	}

	return entities, nil
}

func (uc *usecase) PlaceOrder(ctx context.Context, model *domain.PlaceOrderModel) error {
	order, err := domain.CreateOrderFrom(ctx, model, uc.productDomainSvc)
	if err != nil {
		return errors.Wrap(err, "usecase-domain.CreateOrderFrom")
	}

	err = uc.orderRepo.Create(ctx, order)
	if err != nil {
		return errors.Wrap(err, "usecase-uc.orderRepo.Create")
	}

	slog.Debug("order created", "order", *order)

	baristaEvents := make([]shared.DomainEvent, 0)
	kitchenEvents := make([]shared.DomainEvent, 0)

	for _, event := range order.DomainEvents() {
		if event.Identity() == "BaristaOrdered" {
			baristaEvents = append(baristaEvents, event)
		}

		if event.Identity() == "KitchenOrdered" {
			kitchenEvents = append(kitchenEvents, event)
		}
	}

	err = uc.baristaEventPub.Publish(ctx, baristaEvents)
	if err != nil {
		return errors.Wrap(err, "usecase-baristaEventPub.Publish")
	}

	err = uc.kitchenEventPub.Publish(ctx, kitchenEvents)
	if err != nil {
		return errors.Wrap(err, "usecase-kitchenEventPub.Publish")
	}

	return nil
}
