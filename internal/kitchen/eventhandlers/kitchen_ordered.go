package eventhandlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/kitchen/domain"
	"github.com/thangchung/go-coffeeshop/internal/pkg/event"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
)

type KitchenOrderedEventHandler interface {
	Handle(context.Context, *event.KitchenOrdered) error
}

var _ KitchenOrderedEventHandler = (*kitchenOrderedEventHandler)(nil)

type kitchenOrderedEventHandler struct {
	repo       domain.OrderRepo
	counterPub *publisher.Publisher
	logger     *mylogger.Logger
}

func NewKitchenOrderedEventHandler(
	repo domain.OrderRepo,
	counterPub *publisher.Publisher,
	logger *mylogger.Logger,
) KitchenOrderedEventHandler {
	return &kitchenOrderedEventHandler{
		repo:       repo,
		counterPub: counterPub,
		logger:     logger,
	}
}

func (h *kitchenOrderedEventHandler) Handle(ctx context.Context, e *event.KitchenOrdered) error {
	h.logger.Info("kitchenOrderedEventHandler-Handle: %v", e)

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

func calculateDelay(itemType shared.ItemType) time.Duration {
	switch itemType {
	case shared.ItemTypeCroissant:
		return 7 * time.Second
	case shared.ItemTypeCroissantChocolate:
		return 7 * time.Second
	case shared.ItemTypeCakePop:
		return 5 * time.Second
	case shared.ItemTypeMuffin:
		return 7 * time.Second
	default:
		return 3 * time.Second
	}
}
