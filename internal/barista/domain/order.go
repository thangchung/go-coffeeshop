package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type BaristaOrder struct {
	ID       uuid.UUID    `json:"id" db:"id"`
	ItemName string       `json:"itemName" db:"item_name"`
	ItemType gen.ItemType `json:"itemType" db:"item_type"`
	TimeUp   time.Time    `json:"timeUp" db:"time_up"`
	Created  time.Time    `json:"created" db:"created"`
	Updated  time.Time    `json:"updated" db:"updated"`
}
