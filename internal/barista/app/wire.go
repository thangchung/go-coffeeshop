//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/eventhandlers"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	pkgConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
)

func InitApp(
	cfg *config.Config,
	pg *postgres.Postgres,
	amqpConn *amqp.Connection,
	counterOrderPub rabbitmq.EventPublisher,
	consumer *pkgConsumer.Consumer,
) (*App, error) {
	panic(wire.Build(New, eventhandlers.BaristaOrderedEventHandlerSet))
}
