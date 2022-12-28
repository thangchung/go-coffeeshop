package domain

import (
	"context"
)

type (
	ProductDomainService interface {
		GetItemsByType(context.Context, *PlaceOrderModel, bool) ([]*ItemModel, error)
	}
)
