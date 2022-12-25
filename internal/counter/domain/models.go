package domain

import (
	"time"

	"github.com/google/uuid"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type PlaceOrderModel struct {
	CommandType     shared.CommandType
	OrderSource     shared.OrderSource
	Location        shared.Location
	LoyaltyMemberID uuid.UUID
	BaristaItems    []*OrderItemModel
	KitchenItems    []*OrderItemModel
	Timestamp       time.Time
}

type OrderItemModel struct {
	ItemType shared.ItemType
}

type ItemModel struct {
	ItemType shared.ItemType
	Price    float64
}
