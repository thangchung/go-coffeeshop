package sharedkernel

import (
	"context"
	"time"
)

type DomainEvent interface {
	CreateAt() time.Time
	Identity() string
}

type EventPublisher interface {
	Publish(context.Context, []DomainEvent) error
}

type AggregateRoot struct {
	domainEvents []DomainEvent
}

func (ar *AggregateRoot) ApplyDomain(e DomainEvent) {
	ar.domainEvents = append(ar.domainEvents, e)
}

func (ar *AggregateRoot) DomainEvents() []DomainEvent {
	return ar.domainEvents
}
