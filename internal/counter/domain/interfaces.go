package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type (
	OrderRepo interface {
		GetAll(context.Context) ([]*Order, error)
		GetByID(context.Context, uuid.UUID) (*Order, error)
		Create(context.Context, *gen.OrderDto) error
		Update(context.Context, *gen.OrderDto) (*gen.OrderDto, error)
	}

	ProductDomainService interface {
		GetItemsByType(context.Context, *gen.PlaceOrderRequest, bool) (*gen.GetItemsByTypeResponse, error)
	}

	BaristaOrderUpdatedEventHandler interface {
		Handle(context.Context, *event.BaristaOrderUpdated) error
	}
)
