package entity

import (
	"github.com/google/uuid"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type Order struct {
	ID              uuid.UUID
	OrderSource     gen.OrderSource
	LoyaltyMemberID uuid.UUID
	OrderStatus     gen.Status
	Location        gen.Location
	LineItems       []LineItem
}

func NewOrder(orderSource gen.OrderSource, loyaltyMemberID uuid.UUID, orderStatus gen.Status, location gen.Location) *Order {
	return &Order{
		ID:              uuid.New(),
		OrderSource:     orderSource,
		LoyaltyMemberID: loyaltyMemberID,
		OrderStatus:     orderStatus,
		Location:        location,
	}
}

func (o *Order) From(request *gen.PlaceOrderRequest) (*Order, error) {
	loyaltyMemberID, err := uuid.Parse(request.LoyaltyMemberId)
	if err != nil {
		return nil, err
	}

	order := NewOrder(request.OrderSource, loyaltyMemberID, gen.Status_IN_PROGRESS, request.Location)

	if len(request.BaristaItems) > 0 {

	}

	if len(request.KitchenItems) > 0 {

	}

	return order, nil
}
