package consumer

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/internal/barista/features/orders/eventhandlers"
	"github.com/thangchung/go-coffeeshop/pkg/event"
	log "github.com/thangchung/go-coffeeshop/pkg/logger"
)

const (
	_exchangeKind       = "direct"
	_exchangeDurable    = true
	_exchangeAutoDelete = false
	_exchangeInternal   = false
	_exchangeNoWait     = false

	_queueDurable    = true
	_queueAutoDelete = false
	_queueExclusive  = false
	_queueNoWait     = false

	_prefetchCount  = 1
	_prefetchSize   = 0
	_prefetchGlobal = false

	_consumeAutoAck   = false
	_consumeExclusive = false
	_consumeNoLocal   = false
	_consumeNoWait    = false

	_exchangeName    = "orders-exchange"
	_queueName       = "orders-queue"
	_bindingKey      = "orders-routing-key"
	_consumerTag     = "orders-consumer"
	_messageTypeName = "ordered"
	_workerPoolSize  = 24
)

type Consumer struct {
	exchangeName, queueName, bindingKey, consumerTag string
	messageTypeName                                  string
	workerPoolSize                                   int
	amqpConn                                         *amqp.Connection
	logger                                           *log.Logger
	handler                                          eventhandlers.BaristaOrderedEventHandler
}

func NewConsumer(
	amqpConn *amqp.Connection,
	handler eventhandlers.BaristaOrderedEventHandler,
	logger *log.Logger,
	opts ...Option,
) (*Consumer, error) {
	sub := &Consumer{
		amqpConn:        amqpConn,
		logger:          logger,
		handler:         handler,
		exchangeName:    _exchangeName,
		queueName:       _queueName,
		bindingKey:      _bindingKey,
		consumerTag:     _consumerTag,
		messageTypeName: _messageTypeName,
		workerPoolSize:  _workerPoolSize,
	}

	for _, opt := range opts {
		opt(sub)
	}

	return sub, nil
}

// CreateChannel Consume messages.
func (c *Consumer) CreateChannel() (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	c.logger.Info("Declaring exchange: %s", c.exchangeName)
	err = ch.ExchangeDeclare(
		c.exchangeName,
		_exchangeKind,
		_exchangeDurable,
		_exchangeAutoDelete,
		_exchangeInternal,
		_exchangeNoWait,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		c.queueName,
		_queueDurable,
		_queueAutoDelete,
		_queueExclusive,
		_queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	c.logger.Info("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		c.exchangeName,
		c.bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		_queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	c.logger.Info("Queue bound to exchange, starting to consume from queue, consumerTag: %v", c.consumerTag)

	err = ch.Qos(
		_prefetchCount,  // prefetch count
		_prefetchSize,   // prefetch size
		_prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}

func (c *Consumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		c.logger.Info("processDeliveries deliveryTag% v", delivery.DeliveryTag)

		switch delivery.Type {
		case c.messageTypeName:
			var payload event.BaristaOrdered
			err := json.Unmarshal(delivery.Body, &payload)

			if err != nil {
				c.logger.LogError(err)
			}

			err = c.handler.Handle(ctx, &payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					c.logger.Error("Err delivery.Reject: %v", err)
				}

				c.logger.Error("Failed to process delivery: %v", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					c.logger.Error("Failed to acknowledge delivery: %v", err)
				}
			}
		default:
			c.logger.Info("default")
		}
	}

	c.logger.Info("Deliveries channel closed")
}

// StartConsumer Start new rabbitmq consumer.
func (c *Consumer) StartConsumer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		c.queueName,
		c.consumerTag,
		_consumeAutoAck,
		_consumeExclusive,
		_consumeNoLocal,
		_consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	forever := make(chan bool)

	for i := 0; i < c.workerPoolSize; i++ {
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.logger.Error("ch.NotifyClose: %v", chanErr)
	<-forever

	return chanErr
}
