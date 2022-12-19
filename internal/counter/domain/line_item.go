package domain

import (
	"github.com/google/uuid"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type LineItem struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	ItemType       shared.ItemType `json:"item_type" db:"item_type"`
	Name           string          `json:"name" db:"name"`
	Price          float32         `json:"price" db:"price"`
	ItemStatus     shared.Status   `json:"item_status" db:"item_status"`
	IsBaristaOrder bool            `json:"is_barista_order" db:"is_barista_order"`
	OrderID        uuid.UUID       `json:"order_id" db:"order_id"` // shadow field
}

func NewLineItem(itemType shared.ItemType, name string, price float32, itemStatus shared.Status, isBarista bool) *LineItem {
	return &LineItem{
		ID:             uuid.New(),
		ItemType:       itemType,
		Name:           name,
		Price:          price,
		ItemStatus:     itemStatus,
		IsBaristaOrder: isBarista,
	}
}
