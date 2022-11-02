package domain

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
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

func CreateOrderFrom(request *gen.PlaceOrderRequest, productServiceClient ProductServiceClient) (*Order, error) {
	loyaltyMemberID, err := uuid.Parse(request.LoyaltyMemberId)
	if err != nil {
		return nil, err
	}

	order := NewOrder(request.OrderSource, loyaltyMemberID, gen.Status_IN_PROGRESS, request.Location)

	numberOfBaristaItems := len(request.BaristaItems) > 0
	numberOfKitchenItems := len(request.KitchenItems) > 0

	if numberOfBaristaItems {
		itemTypesRes, err := productServiceClient.GetItemsByType(request, true)
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.BaristaItems, func(item *gen.CommandItem, index int) {
			find, ok := lo.Find(itemTypesRes.Items, func(i *gen.ItemDto) bool {
				return i.Type == int32(item.ItemType)
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), gen.Status_IN_PROGRESS, true)

				//TODO: add domain events
				// ...

				order.LineItems = append(order.LineItems, *lineItem)
			}
		})
	}

	if numberOfKitchenItems {
		itemTypesRes, err := productServiceClient.GetItemsByType(request, false)
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.KitchenItems, func(item *gen.CommandItem, index int) {
			find, ok := lo.Find(itemTypesRes.Items, func(i *gen.ItemDto) bool {
				return i.Type == int32(item.ItemType)
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), gen.Status_IN_PROGRESS, false)

				//TODO: add domain events
				// ...

				order.LineItems = append(order.LineItems, *lineItem)
			}
		})
	}

	return order, nil
}
