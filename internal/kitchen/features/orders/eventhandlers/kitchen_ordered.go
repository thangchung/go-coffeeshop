package eventhandlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/kitchen/domain"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
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
		ItemType: convertToItemType(e.ItemType),
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

func convertToItemType(dto gen.ItemType) domain.ItemType {
	switch dto {
	case gen.ItemType_CROISSANT:
		return domain.Croissant
	case gen.ItemType_CROISSANT_CHOCOLATE:
		return domain.CroissantChocolate
	case gen.ItemType_CAKEPOP:
		return domain.CakePop
	case gen.ItemType_MUFFIN:
		return domain.Muffin
	default:
		return domain.Croissant
	}
}
