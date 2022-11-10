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
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

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

func (p *OrderPublisher) SetupExchangeAndQueue(exchange, queueName, bindingKey, consumerTag string) error {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	p.logger.Info("Declaring exchange: %s", exchange)

	err = p.amqpChan.ExchangeDeclare(
		exchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := p.amqpChan.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Error ch.QueueDeclare")
	}

	p.logger.Info("Declared queue, binding it to exchange: Queue: %v, "+
		"consumerCount: %v, exchange: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchange,
		bindingKey,
	)

	err = p.amqpChan.QueueBind(
		queue.Name,
		bindingKey,
		exchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Error ch.QueueBind")
	}

	p.logger.Info("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag)

	return nil
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
