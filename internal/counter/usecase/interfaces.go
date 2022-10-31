package usecase

import (
	"context"

	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type (
	QueryOrderFulfillmentUseCase interface {
		GetListOrderFulfillment(context.Context) ([]gen.OrderDto, error)
	}

	QueryOrderFulfillmentRepo interface {
		GetListOrderFulfillment(context.Context) ([]gen.OrderDto, error)
	}
)
