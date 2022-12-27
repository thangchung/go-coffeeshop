//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/app/router"
	"github.com/thangchung/go-coffeeshop/internal/counter/events/handlers"
	infrasGRPC "github.com/thangchung/go-coffeeshop/internal/counter/infras/grpc"
	"github.com/thangchung/go-coffeeshop/internal/counter/infras/repo"
	ordersUC "github.com/thangchung/go-coffeeshop/internal/counter/usecases/orders"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	pkgConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"google.golang.org/grpc"
)

func baristaEventPublisher(amqpConn *amqp091.Connection) ordersUC.BaristaEventPublisher {
	pub, _ := pkgPublisher.NewPublisher(amqpConn)
	return (ordersUC.BaristaEventPublisher)(pub)
}

func kitchenEventPublisher(amqpConn *amqp091.Connection) ordersUC.KitchenEventPublisher {
	pub, _ := pkgPublisher.NewPublisher(amqpConn)
	return (ordersUC.KitchenEventPublisher)(pub)
}

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
	grpcServer *grpc.Server,
) (*App, error) {
	panic(wire.Build(
		New,
		postgres.DBEngineSet,
		rabbitmq.RabbitMQSet,
		pkgConsumer.EventConsumerSet,
		infrasGRPC.ProductGRPCClientSet,
		router.CounterGRPCServerSet,
		repo.RepositorySet,
		ordersUC.UseCaseSet,
		baristaEventPublisher,
		kitchenEventPublisher,
		handlers.BaristaOrderUpdatedEventHandlerSet,
		handlers.KitchenOrderUpdatedEventHandlerSet,
	))
}
