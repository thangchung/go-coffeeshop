package sharedkernel

import (
	"context"
	"time"
)

type (
	DomainEvent interface {
		CreateAt() time.Time
		Identity() string
	}

	EventPublisher interface {
		Publish(context.Context, []DomainEvent) error
	}
)
