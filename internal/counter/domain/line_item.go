package domain

import (
	"github.com/google/uuid"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type LineItem struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	ItemType       gen.ItemType `json:"item_type" db:"item_type"`
	Name           string       `json:"name" db:"name"`
	Price          float32      `json:"price" db:"price"`
	ItemStatus     gen.Status   `json:"item_status" db:"item_status"`
	IsBaristaOrder bool         `json:"is_barista_order" db:"is_barista_order"`
	OrderID        uuid.UUID    `json:"order_id" db:"order_id"` // shadow field
}

func NewLineItem(itemType gen.ItemType, name string, price float32, itemStatus gen.Status, isBarista bool) *LineItem {
	return &LineItem{
		ID:             uuid.New(),
		ItemType:       itemType,
		Name:           name,
		Price:          price,
		ItemStatus:     itemStatus,
		IsBaristaOrder: isBarista,
	}
}
