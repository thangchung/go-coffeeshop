package domain

import (
	"time"

	"github.com/google/uuid"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

type BaristaOrder struct {
	ID       uuid.UUID       `json:"id" db:"id"`
	ItemName string          `json:"itemName" db:"item_name"`
	ItemType shared.ItemType `json:"itemType" db:"item_type"`
	TimeUp   time.Time       `json:"timeUp" db:"time_up"`
	Created  time.Time       `json:"created" db:"created"`
	Updated  time.Time       `json:"updated" db:"updated"`
}
