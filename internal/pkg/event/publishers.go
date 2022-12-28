package event

import (
	"context"

	"github.com/google/wire"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
)

var (
	BaristaEventPublisherSet = wire.NewSet(NewBaristaEventPublisher)
	KitchenEventPublisherSet = wire.NewSet(NewKitchenEventPublisher)
)

type BaristaEventPublisher interface {
	Configure(...publisher.Option)
	CloseChan()
	Publish(context.Context, []byte, string) error
}

type KitchenEventPublisher interface {
	Configure(...publisher.Option)
	CloseChan()
	Publish(context.Context, []byte, string) error
}

type baristaEventPublisher struct {
	pub publisher.EventPublisher
}

func NewBaristaEventPublisher(pub publisher.EventPublisher) BaristaEventPublisher {
	return &baristaEventPublisher{
		pub: pub,
	}
}

func (p *baristaEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

func (p *baristaEventPublisher) CloseChan() {
	p.pub.CloseChan()
}

func (p *baristaEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

type kitchenEventPublisher struct {
	pub publisher.EventPublisher
}

func NewKitchenEventPublisher(pub publisher.EventPublisher) KitchenEventPublisher {
	return &kitchenEventPublisher{
		pub: pub,
	}
}

func (p *kitchenEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

func (p *kitchenEventPublisher) CloseChan() {
	p.pub.CloseChan()
}

func (p *kitchenEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}
