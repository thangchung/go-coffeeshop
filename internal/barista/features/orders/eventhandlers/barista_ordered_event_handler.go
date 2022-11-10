package eventhandlers

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/pkg/event"
)

type BaristaOrderedEventHandler interface {
	Handle(context.Context, *event.BaristaOrdered) error
}

type DefaultBaristaOrderedEventHandler struct{}

func NewBaristaOrderedEventHandler() *DefaultBaristaOrderedEventHandler {
	return &DefaultBaristaOrderedEventHandler{}
}

func (h *DefaultBaristaOrderedEventHandler) Handle(ctx context.Context, e *event.BaristaOrdered) error {
	fmt.Println(e)

	return nil
}
