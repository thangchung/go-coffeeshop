package rabbitmq

import "context"

type EventPublisher interface {
	Publish(context.Context, []byte, string) error
	CloseChan()
}
