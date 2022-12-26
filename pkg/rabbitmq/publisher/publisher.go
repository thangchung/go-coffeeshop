package publisher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	"golang.org/x/exp/slog"
)

const (
	_publishMandatory = false
	_publishImmediate = false

	_exchangeName    = "orders-exchange"
	_bindingKey      = "orders-routing-key"
	_messageTypeName = "ordered"
)

type publisher struct {
	exchangeName, bindingKey string
	messageTypeName          string
	amqpChan                 *amqp.Channel
	amqpConn                 *amqp.Connection
}

var _ rabbitmq.EventPublisher = (*publisher)(nil)

func NewPublisher(amqpConn *amqp.Connection, opts ...Option) (rabbitmq.EventPublisher, error) {
	ch, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	pub := &publisher{
		amqpConn:        amqpConn,
		amqpChan:        ch,
		exchangeName:    _exchangeName,
		bindingKey:      _bindingKey,
		messageTypeName: _messageTypeName,
	}

	for _, opt := range opts {
		opt(pub)
	}

	return pub, nil
}

// CloseChan Close messages chan.
func (p *publisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		slog.Error("failed to close chan", err)
	}
}

func (p *publisher) PublishEvents(ctx context.Context, events []any) error {
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return errors.Wrap(err, "publisher-json.Marshal")
		}

		err = p.Publish(ctx, b, "text/plain")
		if err != nil {
			return errors.Wrap(err, "publisher-pub.Publish")
		}
	}

	return nil
}

// Publish message.
func (p *publisher) Publish(ctx context.Context, body []byte, contentType string) error {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	slog.Info("publish message", "exchange", p.exchangeName, "routing_key", p.bindingKey)

	if err := ch.PublishWithContext(
		ctx,
		p.exchangeName,
		p.bindingKey,
		_publishMandatory,
		_publishImmediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
			Type:         p.messageTypeName,
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}

	return nil
}
