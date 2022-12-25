package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	events "github.com/thangchung/go-coffeeshop/internal/pkg/event"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type Order struct {
	shared.AggregateRoot
	ID              uuid.UUID
	OrderSource     shared.OrderSource
	LoyaltyMemberID uuid.UUID
	OrderStatus     shared.Status
	Location        shared.Location
	LineItems       []*LineItem
}

func NewOrder(
	orderSource shared.OrderSource,
	loyaltyMemberID uuid.UUID,
	orderStatus shared.Status,
	location shared.Location,
) *Order {
	return &Order{
		ID:              uuid.New(),
		OrderSource:     orderSource,
		LoyaltyMemberID: loyaltyMemberID,
		OrderStatus:     orderStatus,
		Location:        location,
	}
}

func CreateOrderFrom(
	ctx context.Context,
	request *PlaceOrderModel,
	productDomainSvc ProductDomainService,
) (*Order, error) {
	order := NewOrder(request.OrderSource, request.LoyaltyMemberID, shared.StatusInProcess, request.Location)

	numberOfBaristaItems := len(request.BaristaItems) > 0
	numberOfKitchenItems := len(request.KitchenItems) > 0

	if numberOfBaristaItems {
		itemTypesRes, err := productDomainSvc.GetItemsByType(ctx, request, true)
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.BaristaItems, func(item *OrderItemModel, _ int) {
			find, ok := lo.Find(itemTypesRes, func(i *ItemModel) bool {
				return i.ItemType == item.ItemType
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), shared.StatusInProcess, true)

				event := events.BaristaOrdered{
					OrderID:    order.ID,
					ItemLineID: lineItem.ID,
					ItemType:   item.ItemType,
				}

				order.ApplyDomain(event)

				order.LineItems = append(order.LineItems, lineItem)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	if numberOfKitchenItems {
		itemTypesRes, err := productDomainSvc.GetItemsByType(ctx, request, false)
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.KitchenItems, func(item *OrderItemModel, index int) {
			find, ok := lo.Find(itemTypesRes, func(i *ItemModel) bool {
				return i.ItemType == item.ItemType
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), shared.StatusInProcess, false)

				event := events.KitchenOrdered{
					OrderID:    order.ID,
					ItemLineID: lineItem.ID,
					ItemType:   item.ItemType,
				}

				order.ApplyDomain(event)

				order.LineItems = append(order.LineItems, lineItem)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (o *Order) Apply(event *events.OrderUp) error {
	if len(o.LineItems) == 0 {
		return nil // we dont do anything
	}

	_, index, ok := lo.FindIndexOf(o.LineItems, func(i *LineItem) bool {
		return i.ItemType == event.ItemType
	})

	if !ok {
		return ErrItemNotFound
	}

	o.LineItems[index].ItemStatus = shared.StatusFulfilled

	if checkFulfilledStatus(o.LineItems) {
		o.OrderStatus = shared.StatusFulfilled
	}

	return nil
}

func checkFulfilledStatus(lineItems []*LineItem) bool {
	for _, item := range lineItems {
		if item.ItemStatus != shared.StatusFulfilled {
			return false
		}
	}

	return true
}
