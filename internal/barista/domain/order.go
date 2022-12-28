package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type BaristaOrder struct {
	shared.AggregateRoot
	ID       uuid.UUID
	ItemName string
	ItemType shared.ItemType
	TimeUp   time.Time
	Created  time.Time
	Updated  time.Time
}

func NewBaristaOrder(e event.BaristaOrdered) BaristaOrder {
	timeIn := time.Now()

	delay := calculateDelay(e.ItemType)
	time.Sleep(delay) // simulate the delay when makes the drink

	timeUp := time.Now().Add(delay)

	order := BaristaOrder{
		ID:       e.ItemLineID,
		ItemName: e.ItemType.String(),
		ItemType: e.ItemType,
		TimeUp:   timeUp,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	orderUpdatedEvent := event.BaristaOrderUpdated{
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
	case shared.ItemTypeCoffeeBlack:
		return 5 * time.Second
	case shared.ItemTypeCoffeeWithRoom:
		return 5 * time.Second
	case shared.ItemTypeEspresso:
		return 7 * time.Second
	case shared.ItemTypeEspressoDouble:
		return 7 * time.Second
	case shared.ItemTypeCappuccino:
		return 10 * time.Second
	default:
		return 3 * time.Second
	}
}
