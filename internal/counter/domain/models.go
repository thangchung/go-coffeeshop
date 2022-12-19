package domain

import (
	"time"

	"github.com/google/uuid"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type OrderModel struct {
	ID              uuid.UUID
	OrderSource     shared.OrderSource
	LoyaltyMemberID uuid.UUID
	OrderStatus     shared.Status
	Location        shared.Location
	LineItems       []*LineItemModel
}

type LineItemModel struct {
	ID             uuid.UUID
	ItemType       shared.ItemType
	Name           string
	Price          float64
	ItemStatus     shared.Status
	IsBaristaOrder bool
}

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
