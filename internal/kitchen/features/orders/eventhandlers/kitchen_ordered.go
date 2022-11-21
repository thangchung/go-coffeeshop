package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/kitchen/domain"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type KitchenOrderedEventHandler interface {
	Handle(context.Context, *event.KitchenOrdered) error
}

var _ KitchenOrderedEventHandler = (*kitchenOrderedEventHandler)(nil)

type kitchenOrderedEventHandler struct {
	repo       domain.OrderRepo
	counterPub *publisher.Publisher
}

func NewKitchenOrderedEventHandler(repo domain.OrderRepo, counterPub *publisher.Publisher) KitchenOrderedEventHandler {
	return &kitchenOrderedEventHandler{
		repo:       repo,
		counterPub: counterPub,
	}
}

func (h *kitchenOrderedEventHandler) Handle(ctx context.Context, e *event.KitchenOrdered) error {
	fmt.Println(e)

	timeIn := time.Now()

	delay := calculateDelay(e.ItemType)
	time.Sleep(delay)

	timeUp := time.Now().Add(delay)

	err := h.repo.Create(ctx, &domain.KitchenOrder{
		ID:       e.ItemLineID,
		OrderID:  e.OrderID,
		ItemType: e.ItemType,
		ItemName: e.ItemType.String(),
		TimeUp:   timeUp,
		Created:  time.Now(),
		Updated:  time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "kitchenOrderedEventHandler-h.repo.Create")
	}

	message := event.KitchenOrderUpdated{
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
