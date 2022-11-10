package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/internal/barista/features/orders/eventhandlers"
	"github.com/thangchung/go-coffeeshop/pkg/event"
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

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type OrderConsumer struct {
	amqpConn *amqp.Connection
	logger   *log.Logger
	handler  eventhandlers.BaristaOrderedEventHandler
}

func NewOrderConsumer(amqpConn *amqp.Connection, handler eventhandlers.BaristaOrderedEventHandler, logger *log.Logger) (*OrderConsumer, error) {
	// ch, err := amqpConn.Channel()
	// if err != nil {
	// 	panic(err)
	// }
	// defer ch.Close()

	return &OrderConsumer{
		amqpConn: amqpConn,
		logger:   logger,
		handler:  handler,
	}, nil
}

// CreateChannel Consume messages
func (c *OrderConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	c.logger.Info("Declaring exchange: %s", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
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
		exchangeName,
		bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	c.logger.Info("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag)

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}

func (c *OrderConsumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		c.logger.Info("processDeliveries deliveryTag% v", delivery.DeliveryTag)

		switch delivery.Type {
		case "barista.ordered":
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
			fmt.Println("default")
		}
	}

	c.logger.Info("Deliveries channel closed")
}

// StartConsumer Start new rabbitmq consumer
func (c *OrderConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	forever := make(chan bool)

	for i := 0; i < workerPoolSize; i++ {
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.logger.Error("ch.NotifyClose: %v", chanErr)
	<-forever

	return chanErr
}
