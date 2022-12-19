package orders

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
)

type baristaOrderUpdatedEventHandler struct {
	orderRepo domain.OrderRepo
}

var _ BaristaOrderUpdatedEventHandler = (*baristaOrderUpdatedEventHandler)(nil)

func NewBaristaOrderUpdatedEventHandler(orderRepo domain.OrderRepo) BaristaOrderUpdatedEventHandler {
	return &baristaOrderUpdatedEventHandler{
		orderRepo: orderRepo,
	}
}

func (h *baristaOrderUpdatedEventHandler) Handle(ctx context.Context, e *event.BaristaOrderUpdated) error {
	order, err := h.orderRepo.GetByID(ctx, e.OrderID)
	if err != nil {
		return fmt.Errorf("NewBaristaOrderUpdatedEventHandler-Handle-h.orderRepo.GetOrderByID(ctx, e.OrderID): %w", err)
	}

	if err = order.Apply(e); err != nil {
		return fmt.Errorf("NewBaristaOrderUpdatedEventHandler-Handle-order.Apply(e): %w", err)
	}

	_, err = h.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("NewBaristaOrderUpdatedEventHandler-Handle-h.orderRepo.Update(ctx, ToDto(order)): %w", err)
	}

	return nil
}
