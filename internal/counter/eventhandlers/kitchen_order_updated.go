package eventhandlers

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
)

type kitchenOrderUpdatedEventHandler struct {
	orderRepo domain.OrderRepo
}

var _ domain.KitchenOrderUpdatedEventHandler = (*kitchenOrderUpdatedEventHandler)(nil)

func NewKitchenOrderUpdatedEventHandler(orderRepo domain.OrderRepo) domain.KitchenOrderUpdatedEventHandler {
	return &kitchenOrderUpdatedEventHandler{
		orderRepo: orderRepo,
	}
}

func (h *kitchenOrderUpdatedEventHandler) Handle(ctx context.Context, e *event.KitchenOrderUpdated) error {
	order, err := h.orderRepo.GetByID(ctx, e.OrderID)
	if err != nil {
		return fmt.Errorf("NewKitchenOrderUpdatedEventHandler-Handle-h.orderRepo.GetOrderByID(ctx, e.OrderID): %w", err)
	}

	orderUp := event.OrderUp{
		OrderID:    e.OrderID,
		ItemLineID: e.ItemLineID,
		Name:       e.Name,
		ItemType:   e.ItemType,
		TimeUp:     e.TimeUp,
		MadeBy:     e.MadeBy,
	}

	if err = order.Apply(&orderUp); err != nil {
		return fmt.Errorf("NewKitchenOrderUpdatedEventHandler-Handle-order.Apply(e): %w", err)
	}

	_, err = h.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("NewKitchenOrderUpdatedEventHandler-Handle-h.orderRepo.Update(ctx, ToDto(order)): %w", err)
	}

	return nil
}
