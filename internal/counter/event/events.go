package event

import (
	"time"

	"github.com/google/uuid"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type BaristaOrdered struct {
	OrderID    uuid.UUID    `json:"orderId"`
	ItemLineID uuid.UUID    `json:"itemLineId"`
	ItemType   gen.ItemType `json:"itemType"`
}

type KitchenOrdered struct {
	OrderID    uuid.UUID    `json:"orderId"`
	ItemLineID uuid.UUID    `json:"itemLineId"`
	ItemType   gen.ItemType `json:"itemType"`
}

type BaristaOrderUpdated struct {
	OrderID    uuid.UUID    `json:"orderId"`
	ItemLineID uuid.UUID    `json:"itemLineId"`
	Name       string       `json:"name"`
	ItemType   gen.ItemType `json:"itemType"`
	TimeIn     time.Time    `json:"timeIn"`
	MadeBy     string       `json:"madeBy"`
	TimeUp     time.Time    `json:"timeUp"`
}

type KitchenOrderUpdated struct {
	OrderID    uuid.UUID    `json:"orderId"`
	ItemLineID uuid.UUID    `json:"itemLineId"`
	Name       string       `json:"name"`
	ItemType   gen.ItemType `json:"itemType"`
	TimeIn     time.Time    `json:"timeIn"`
	MadeBy     string       `json:"madeBy"`
	TimeUp     time.Time    `json:"timeUp"`
}
