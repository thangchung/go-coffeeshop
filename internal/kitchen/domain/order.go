package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type KitchenOrder struct {
	shared.AggregateRoot
	ID       uuid.UUID
	OrderID  uuid.UUID
	ItemName string
	ItemType shared.ItemType
	TimeUp   time.Time
	Created  time.Time
	Updated  time.Time
}

func NewKitchenOrder(e event.KitchenOrdered) KitchenOrder {
	timeIn := time.Now()

	delay := calculateDelay(e.ItemType)
	time.Sleep(delay) // simulate the delay when makes the drink

	timeUp := time.Now().Add(delay)

	order := KitchenOrder{
		ID:       e.ItemLineID,
		OrderID:  e.OrderID,
		ItemName: e.ItemType.String(),
		ItemType: e.ItemType,
		TimeUp:   timeUp,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	orderUpdatedEvent := event.KitchenOrderUpdated{
		OrderID:    e.OrderID,
		ItemLineID: e.ItemLineID,
		Name:       e.ItemType.String(),
		ItemType:   e.ItemType,
		MadeBy:     "teesee",
		TimeIn:     timeIn,
		TimeUp:     timeUp,
	}

	order.ApplyDomain(&orderUpdatedEvent)

	return order
}

func calculateDelay(itemType shared.ItemType) time.Duration {
	switch itemType {
	case shared.ItemTypeCroissant:
		return 7 * time.Second
	case shared.ItemTypeCroissantChocolate:
		return 7 * time.Second
	case shared.ItemTypeCakePop:
		return 5 * time.Second
	case shared.ItemTypeMuffin:
		return 7 * time.Second
	default:
		return 3 * time.Second
	}
}
