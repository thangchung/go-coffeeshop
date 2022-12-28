package orders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	"golang.org/x/exp/slog"
)

type usecase struct {
	orderRepo        domain.OrderRepo
	productDomainSvc domain.ProductDomainService
	baristaEventPub  event.BaristaEventPublisher
	kitchenEventPub  event.KitchenEventPublisher
}

var _ UseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	orderRepo domain.OrderRepo,
	productDomainSvc domain.ProductDomainService,
	baristaEventPub event.BaristaEventPublisher,
	kitchenEventPub event.KitchenEventPublisher,
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

	// todo: it might cause dual-write problem, but we accept it temporary
	for _, event := range order.DomainEvents() {
		if event.Identity() == "BaristaOrdered" {
			eventBytes, err := json.Marshal(event)
			if err != nil {
				return errors.Wrap(err, "json.Marshal[event]")
			}

			uc.baristaEventPub.Publish(ctx, eventBytes, "text/plain")
		}

		if event.Identity() == "KitchenOrdered" {
			eventBytes, err := json.Marshal(event)
			if err != nil {
				return errors.Wrap(err, "json.Marshal[event]")
			}

			uc.kitchenEventPub.Publish(ctx, eventBytes, "text/plain")
		}
	}

	return nil
}
