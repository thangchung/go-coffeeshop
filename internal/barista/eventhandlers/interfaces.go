package eventhandlers

import (
	"context"

	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
)

type BaristaOrderedEventHandler interface {
	Handle(context.Context, event.BaristaOrdered) error
}
