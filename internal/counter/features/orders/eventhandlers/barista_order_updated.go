package eventhandlers

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type baristaOrderUpdatedEventHandler struct {
	orderRepo domain.OrderRepo
}

var _ domain.BaristaOrderUpdatedEventHandler = (*baristaOrderUpdatedEventHandler)(nil)

func NewBaristaOrderUpdatedEventHandler(orderRepo domain.OrderRepo) domain.BaristaOrderUpdatedEventHandler {
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

	_, err = h.orderRepo.Update(ctx, ToDto(order))
	if err != nil {
		return fmt.Errorf("NewBaristaOrderUpdatedEventHandler-Handle-h.orderRepo.Update(ctx, ToDto(order)): %w", err)
	}

	return nil
}

func ToDto(order *domain.Order) *gen.OrderDto {
	orderModel := &gen.OrderDto{
		Id:              order.ID.String(),
		Localtion:       order.Location,
		OrderSource:     order.OrderSource,
		OrderStatus:     order.OrderStatus,
		LoyaltyMemberId: order.LoyaltyMemberID.String(),
	}

	for _, item := range order.LineItems {
		orderModel.LineItems = append(orderModel.LineItems, &gen.LineItemDto{
			Id:             item.ID.String(),
			Name:           item.Name,
			Price:          float64(item.Price),
			ItemType:       item.ItemType,
			ItemStatus:     item.ItemStatus,
			IsBaristaOrder: item.IsBaristaOrder,
		})
	}

	return orderModel
}
