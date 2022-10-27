package entity

import (
	"github.com/google/uuid"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type LineItem struct {
	ID             uuid.UUID
	ItemType       gen.ItemType
	Name           string
	Price          float32
	ItemStatus     gen.Status
	IsBaristaOrder bool
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
