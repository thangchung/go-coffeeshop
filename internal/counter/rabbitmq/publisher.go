package rabbitmq

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	log "github.com/thangchung/go-coffeeshop/pkg/logger"
)

const (
	publishMandatory = false
	publishImmediate = false
)

type OrderPublisher struct {
	amqpChan *amqp.Channel
	amqpConn *amqp.Connection
	cfg      *config.Config
	logger   *log.Logger
}

func NewOrderPublisher(amqpConn *amqp.Connection, cfg *config.Config, logger *log.Logger) (*OrderPublisher, error) {
	ch, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	return &OrderPublisher{
		amqpConn: amqpConn,
		amqpChan: ch,
		logger:   logger,
		cfg:      cfg,
	}, nil
}

// CloseChan Close messages chan.
func (p *OrderPublisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		p.logger.Error("OrderPublisher CloseChan: %v", err)
	}
}

// Publish message.
func (p *OrderPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	p.logger.Info("Publishing message Exchange: %s, RoutingKey: %s", p.cfg.RabbitMQ.Exchange, p.cfg.RabbitMQ.RoutingKey)

	if err := ch.PublishWithContext(
		ctx,
		p.cfg.RabbitMQ.Exchange,
		p.cfg.RabbitMQ.RoutingKey,
		publishMandatory,
		publishImmediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
			Type:         "barista.ordered", //todo
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}

	return nil
}
