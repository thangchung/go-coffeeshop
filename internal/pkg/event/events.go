package event

import (
	"time"

	"github.com/google/uuid"

	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type BaristaOrdered struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	ItemType   shared.ItemType `json:"itemType"`
}

func (e BaristaOrdered) Identity() string {
	return "BaristaOrdered"
}

type KitchenOrdered struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	ItemType   shared.ItemType `json:"itemType"`
}

func (e KitchenOrdered) Identity() string {
	return "KitchenOrdered"
}

type BaristaOrderUpdated struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	Name       string          `json:"name"`
	ItemType   shared.ItemType `json:"itemType"`
	TimeIn     time.Time       `json:"timeIn"`
	MadeBy     string          `json:"madeBy"`
	TimeUp     time.Time       `json:"timeUp"`
}

func (e *BaristaOrderUpdated) Identity() string {
	return "BaristaOrderUpdated"
}

type KitchenOrderUpdated struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	Name       string          `json:"name"`
	ItemType   shared.ItemType `json:"itemType"`
	TimeIn     time.Time       `json:"timeIn"`
	MadeBy     string          `json:"madeBy"`
	TimeUp     time.Time       `json:"timeUp"`
}

func (e *KitchenOrderUpdated) Identity() string {
	return "KitchenOrderUpdated"
}

type OrderUp struct {
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	Name       string          `json:"name"`
	ItemType   shared.ItemType `json:"itemType"`
	TimeUp     time.Time       `json:"timeUp"`
	MadeBy     string          `json:"madeBy"`
}

func (e *OrderUp) Identity() string {
	return "OrderUp"
}
