//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/eventhandlers"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	pkgConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
) (*App, error) {
	panic(wire.Build(New, postgres.DBEngineSet, rabbitmq.RabbitMQSet, pkgPublisher.EventPublisherSet, pkgConsumer.EventConsumerSet, eventhandlers.BaristaOrderedEventHandlerSet))
}
