package event

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
)

type eventPublisherRabbitMQ struct {
	pub publisher.Publisher
}

var _ shared.EventPublisher = (*eventPublisherRabbitMQ)(nil)

func NewEventPublisher(pub publisher.Publisher) shared.EventPublisher {
	return &eventPublisherRabbitMQ{
		pub: pub,
	}
}

func (p *eventPublisherRabbitMQ) Publish(ctx context.Context, events []shared.DomainEvent) error {
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return errors.Wrap(err, "eventPublisherRabbitMQ-json.Marshal")
		}

		err = p.pub.Publish(ctx, b, "text/plain")
		if err != nil {
			return errors.Wrap(err, "eventPublisherRabbitMQ-pub.Publish")
		}
	}

	return nil
}
