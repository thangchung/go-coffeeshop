package domain

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	events "github.com/thangchung/go-coffeeshop/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type Order struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	OrderSource     gen.OrderSource `json:"order_source" db:"order_source"`
	LoyaltyMemberID uuid.UUID       `json:"loyalty_member_id" db:"loyalty_member_id"`
	OrderStatus     gen.Status      `json:"order_status" db:"order_status"`
	Location        gen.Location    `json:"location" db:"location"`
	LineItems       []*LineItem
}

func NewOrder(
	orderSource gen.OrderSource,
	loyaltyMemberID uuid.UUID,
	orderStatus gen.Status,
	location gen.Location,
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
	request *gen.PlaceOrderRequest,
	productDomainSvc ProductDomainService,
	baristaOrderPub publisher.Publisher,
	kitchenOrderPub publisher.Publisher,
) (*Order, error) {
	loyaltyMemberID, err := uuid.Parse(request.LoyaltyMemberId)
	if err != nil {
		return nil, err
	}

	order := NewOrder(request.OrderSource, loyaltyMemberID, gen.Status_IN_PROGRESS, request.Location)

	numberOfBaristaItems := len(request.BaristaItems) > 0
	numberOfKitchenItems := len(request.KitchenItems) > 0

	if numberOfBaristaItems {
		itemTypesRes, err := productDomainSvc.GetItemsByType(ctx, request, true)
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.BaristaItems, func(item *gen.CommandItem, index int) {
			find, ok := lo.Find(itemTypesRes.Items, func(i *gen.ItemDto) bool {
				return i.Type == int32(item.ItemType)
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), gen.Status_IN_PROGRESS, true)

				err = publishBaristaOrderEvent(
					ctx,
					order.ID,
					lineItem.ID,
					lineItem.ItemType,
					baristaOrderPub,
					kitchenOrderPub,
					true,
				)

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

		lo.ForEach(request.KitchenItems, func(item *gen.CommandItem, index int) {
			find, ok := lo.Find(itemTypesRes.Items, func(i *gen.ItemDto) bool {
				return i.Type == int32(item.ItemType)
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), gen.Status_IN_PROGRESS, false)

				err = publishBaristaOrderEvent(
					ctx,
					order.ID,
					lineItem.ID,
					lineItem.ItemType,
					baristaOrderPub,
					kitchenOrderPub,
					false,
				)

				order.LineItems = append(order.LineItems, lineItem)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (o *Order) Apply(event *events.BaristaOrderUpdated) error {
	if len(o.LineItems) == 0 {
		return nil // we dont do anything
	}

	_, index, ok := lo.FindIndexOf(o.LineItems, func(i *LineItem) bool {
		return i.ItemType == event.ItemType
	})

	if !ok {
		return errors.New("item not found")
	}

	o.LineItems[index].ItemStatus = gen.Status_FULFILLED

	if checkFulfilledStatus(o.LineItems) {
		o.OrderStatus = gen.Status_FULFILLED
	}

	return nil
}

func publishBaristaOrderEvent(
	ctx context.Context,
	orderID uuid.UUID,
	lineItemID uuid.UUID,
	itemType gen.ItemType,
	baristaOrderPub publisher.Publisher,
	kitchenOrderPub publisher.Publisher,
	isBarista bool,
) error {
	if isBarista {
		// todo: refactor to event domain dispatcher
		event := events.BaristaOrdered{
			OrderID:    orderID,
			ItemLineID: lineItemID,
			ItemType:   itemType,
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			return errors.Wrap(err, "json.Marshal - events.BaristaOrdered")
		}

		err = baristaOrderPub.Publish(ctx, eventBytes, "text/plain")
		if err != nil {
			return errors.Wrap(err, "orderPublisher - Publish")
		}

		return nil
	} else {
		// todo: refactor to event domain dispatcher
		event := events.KitchenOrdered{
			OrderID:    orderID,
			ItemLineID: lineItemID,
			ItemType:   itemType,
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			return errors.Wrap(err, "json.Marshal - events.BaristaOrdered")
		}

		err = kitchenOrderPub.Publish(ctx, eventBytes, "text/plain")
		if err != nil {
			return errors.Wrap(err, "orderPublisher - Publish")
		}

		return nil
	}
}

func checkFulfilledStatus(lineItems []*LineItem) bool {
	for _, item := range lineItems {
		if item.ItemStatus != gen.Status_FULFILLED {
			return false
		}
	}

	return true
}
