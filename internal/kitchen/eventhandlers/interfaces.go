package eventhandlers

import (
	"context"

	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
)

type KitchenOrderedEventHandler interface {
	Handle(context.Context, event.KitchenOrdered) error
}
