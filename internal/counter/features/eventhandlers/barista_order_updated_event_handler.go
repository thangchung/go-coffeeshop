package eventhandlers

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/pkg/event"
)

type BaristaOrderUpdatedEventHandler interface {
	Handle(context.Context, *event.BaristaOrderUpdated) error
}

type DefaultBaristaOrderUpdatedEventHandler struct{}

func NewBaristaOrderUpdatedEventHandler() *DefaultBaristaOrderUpdatedEventHandler {
	return &DefaultBaristaOrderUpdatedEventHandler{}
}

func (h *DefaultBaristaOrderUpdatedEventHandler) Handle(ctx context.Context, e *event.BaristaOrderUpdated) error {
	fmt.Println(e)

	return nil
}
