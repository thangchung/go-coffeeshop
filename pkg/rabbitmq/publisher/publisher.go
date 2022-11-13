package publisher

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/thangchung/go-coffeeshop/pkg/logger"
)

const (
	_publishMandatory = false
	_publishImmediate = false

	_exchangeName    = "orders-exchange"
	_bindingKey      = "orders-routing-key"
	_messageTypeName = "ordered"
)

type Publisher struct {
	exchangeName, bindingKey string
	messageTypeName          string
	amqpChan                 *amqp.Channel
	amqpConn                 *amqp.Connection
	logger                   *log.Logger
}

func NewPublisher(amqpConn *amqp.Connection, logger *log.Logger, opts ...Option) (*Publisher, error) {
	ch, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	pub := &Publisher{
		amqpConn:        amqpConn,
		amqpChan:        ch,
		logger:          logger,
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
func (p *Publisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		p.logger.Error("Publisher CloseChan: %v", err)
	}
}

// Publish message.
func (p *Publisher) Publish(ctx context.Context, body []byte, contentType string) error {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	p.logger.Info("Publishing message Exchange: %s, RoutingKey: %s", p.exchangeName, p.bindingKey)

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
			Type:         p.messageTypeName, //"barista.ordered",
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}

	return nil
}
