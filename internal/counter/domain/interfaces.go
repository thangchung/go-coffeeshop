package domain

import (
	"context"

	"github.com/google/uuid"
)

type (
	OrderRepo interface {
		GetAll(context.Context) ([]*Order, error)
		GetByID(context.Context, uuid.UUID) (*Order, error)
		Create(context.Context, *Order) error
		Update(context.Context, *Order) (*Order, error)
	}

	ProductDomainService interface {
		GetItemsByType(context.Context, *PlaceOrderModel, bool) ([]*ItemModel, error)
	}
)
