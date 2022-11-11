package domain

import (
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type (
	QueryOrderFulfillmentRepo interface {
		GetListOrderFulfillment() ([]gen.OrderDto, error)
	}

	ProductDomainService interface {
		GetItemsByType(*gen.PlaceOrderRequest, bool) (*gen.GetItemsByTypeResponse, error)
	}
)
