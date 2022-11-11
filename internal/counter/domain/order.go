package domain

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	counterRabbitMQ "github.com/thangchung/go-coffeeshop/internal/counter/rabbitmq"
	events "github.com/thangchung/go-coffeeshop/pkg/event"
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
	productDomainService ProductDomainService,
	orderPublisher counterRabbitMQ.OrderPublisher,
) (*Order, error) {
	loyaltyMemberID, err := uuid.Parse(request.LoyaltyMemberId)
	if err != nil {
		return nil, err
	}

	order := NewOrder(request.OrderSource, loyaltyMemberID, gen.Status_IN_PROGRESS, request.Location)

	numberOfBaristaItems := len(request.BaristaItems) > 0
	numberOfKitchenItems := len(request.KitchenItems) > 0

	if numberOfBaristaItems {
		itemTypesRes, err := productDomainService.GetItemsByType(request, true)
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.BaristaItems, func(item *gen.CommandItem, index int) {
			find, ok := lo.Find(itemTypesRes.Items, func(i *gen.ItemDto) bool {
				return i.Type == int32(item.ItemType)
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), gen.Status_IN_PROGRESS, true)

				err = publishBaristaOrderEvent(ctx, order.ID, lineItem.ID, lineItem.ItemType, orderPublisher, true)

				order.LineItems = append(order.LineItems, *lineItem)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	if numberOfKitchenItems {
		itemTypesRes, err := productDomainService.GetItemsByType(request, false)
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.KitchenItems, func(item *gen.CommandItem, index int) {
			find, ok := lo.Find(itemTypesRes.Items, func(i *gen.ItemDto) bool {
				return i.Type == int32(item.ItemType)
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), float32(find.Price), gen.Status_IN_PROGRESS, false)

				err = publishBaristaOrderEvent(ctx, order.ID, lineItem.ID, lineItem.ItemType, orderPublisher, false)

				order.LineItems = append(order.LineItems, *lineItem)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func publishBaristaOrderEvent(
	ctx context.Context,
	orderID uuid.UUID,
	lineItemID uuid.UUID,
	itemType gen.ItemType,
	publisher counterRabbitMQ.OrderPublisher,
	isBarista bool,
) error {
	if isBarista {
		// todo: refactor to event domain dispatcher
		// ...
		event := events.BaristaOrdered{
			OrderID:    orderID,
			ItemLineID: lineItemID,
			ItemType:   itemType,
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			return errors.Wrap(err, "json.Marshal - events.BaristaOrdered")
		}

		err = publisher.Publish(ctx, eventBytes, "text/plain")
		if err != nil {
			return errors.Wrap(err, "orderPublisher - Publish")
		}

		return nil
	} else {
		// todo: refactor to event domain dispatcher
		// ...
		event := events.KitchenOrdered{
			OrderID:    orderID,
			ItemLineID: lineItemID,
			ItemType:   itemType,
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			return errors.Wrap(err, "json.Marshal - events.BaristaOrdered")
		}

		err = publisher.Publish(ctx, eventBytes, "text/plain")
		if err != nil {
			return errors.Wrap(err, "orderPublisher - Publish")
		}

		return nil
	}
}
