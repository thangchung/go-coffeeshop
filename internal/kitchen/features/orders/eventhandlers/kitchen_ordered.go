package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type KitchenOrderedEventHandler interface {
	Handle(context.Context, *event.KitchenOrdered) error
}

var _ KitchenOrderedEventHandler = (*kitchenOrderedEventHandler)(nil)

type kitchenOrderedEventHandler struct {
	counterPub *publisher.Publisher
}

func NewKitchenOrderedEventHandler(counterPub *publisher.Publisher) KitchenOrderedEventHandler {
	return &kitchenOrderedEventHandler{
		counterPub: counterPub,
	}
}

func (h *kitchenOrderedEventHandler) Handle(ctx context.Context, e *event.KitchenOrdered) error {
	fmt.Println(e)

	delay := calculateDelay(e.ItemType)
	time.Sleep(delay)

	// todo: save to db
	// ...

	message := event.KitchenOrderUpdated{
		OrderID:    e.OrderID,
		ItemLineID: e.ItemLineID,
		Name:       e.ItemType.String(),
		ItemType:   e.ItemType,
		MadeBy:     "teesee",
		TimeIn:     time.Now(),
		TimeUp:     time.Now().Add(5 * time.Minute),
	}

	eventBytes, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "json.Marshal-events.KitchenOrderUpdated")
	}

	if err := h.counterPub.Publish(ctx, eventBytes, "text/plain"); err != nil {
		return errors.Wrap(err, "KitchenOrderedEventHandler-Publish")
	}

	return nil
}

func calculateDelay(itemType gen.ItemType) time.Duration {
	switch itemType {
	case gen.ItemType_CROISSANT:
		return 7 * time.Second
	case gen.ItemType_CROISSANT_CHOCOLATE:
		return 7 * time.Second
	case gen.ItemType_CAKEPOP:
		return 5 * time.Second
	case gen.ItemType_MUFFIN:
		return 7 * time.Second
	default:
		return 3 * time.Second
	}
}
