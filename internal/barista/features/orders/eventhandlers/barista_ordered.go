package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/barista/domain"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type BaristaOrderedEventHandler interface {
	Handle(context.Context, *event.BaristaOrdered) error
}

var _ BaristaOrderedEventHandler = (*baristaOrderedEventHandler)(nil)

type baristaOrderedEventHandler struct {
	repo       domain.OrderRepo
	counterPub *publisher.Publisher
}

func NewBaristaOrderedEventHandler(repo domain.OrderRepo, counterPub *publisher.Publisher) BaristaOrderedEventHandler {
	return &baristaOrderedEventHandler{
		repo:       repo,
		counterPub: counterPub,
	}
}

func (h *baristaOrderedEventHandler) Handle(ctx context.Context, e *event.BaristaOrdered) error {
	fmt.Println(e)

	timeIn := time.Now()

	delay := calculateDelay(e.ItemType)
	time.Sleep(delay)

	timeUp := time.Now().Add(delay)

	err := h.repo.Create(ctx, &domain.BaristaOrder{
		ID:       e.ItemLineID,
		ItemType: e.ItemType,
		ItemName: e.ItemType.String(),
		TimeUp:   timeUp,
		Created:  time.Now(),
		Updated:  time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "baristaOrderedEventHandler-h.repo.Create")
	}

	message := event.BaristaOrderUpdated{
		OrderID:    e.OrderID,
		ItemLineID: e.ItemLineID,
		Name:       e.ItemType.String(),
		ItemType:   e.ItemType,
		MadeBy:     "teesee",
		TimeIn:     timeIn,
		TimeUp:     timeUp,
	}

	eventBytes, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "json.Marshal - events.BaristaOrderUpdated")
	}

	if err := h.counterPub.Publish(ctx, eventBytes, "text/plain"); err != nil {
		return errors.Wrap(err, "BaristaOrderedEventHandler - Publish")
	}

	return nil
}

func calculateDelay(itemType gen.ItemType) time.Duration {
	switch itemType {
	case gen.ItemType_COFFEE_BLACK:
		return 5 * time.Second
	case gen.ItemType_COFFEE_WITH_ROOM:
		return 5 * time.Second
	case gen.ItemType_ESPRESSO:
		return 7 * time.Second
	case gen.ItemType_ESPRESSO_DOUBLE:
		return 7 * time.Second
	case gen.ItemType_CAPPUCCINO:
		return 10 * time.Second
	default:
		return 3 * time.Second
	}
}
