package orders

import (
	"context"

	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
)

type (
	UseCase interface {
		GetListOrderFulfillment(context.Context) ([]*domain.Order, error)
		PlaceOrder(context.Context, *domain.PlaceOrderModel) error
	}
)
