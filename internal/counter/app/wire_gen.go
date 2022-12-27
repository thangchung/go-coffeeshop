// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/app/router"
	"github.com/thangchung/go-coffeeshop/internal/counter/events/handlers"
	grpc2 "github.com/thangchung/go-coffeeshop/internal/counter/infras/grpc"
	"github.com/thangchung/go-coffeeshop/internal/counter/infras/repo"
	"github.com/thangchung/go-coffeeshop/internal/counter/usecases/orders"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"
	"google.golang.org/grpc"
)

// Injectors from wire.go:

func InitApp(cfg *config.Config, dbConnStr postgres.DBConnString, rabbitMQConnStr rabbitmq.RabbitMQConnStr, grpcServer *grpc.Server) (*App, error) {
	dbEngine, err := postgres.NewPostgresDB(dbConnStr)
	if err != nil {
		return nil, err
	}
	connection, err := rabbitmq.NewRabbitMQConn(rabbitMQConnStr)
	if err != nil {
		return nil, err
	}
	ordersBaristaEventPublisher := baristaEventPublisher(connection)
	ordersKitchenEventPublisher := kitchenEventPublisher(connection)
	eventConsumer, err := consumer.NewConsumer(connection)
	if err != nil {
		return nil, err
	}
	productDomainService, err := grpc2.NewGRPCProductClient(cfg)
	if err != nil {
		return nil, err
	}
	orderRepo := repo.NewOrderRepo(dbEngine)
	useCase := orders.NewUseCase(orderRepo, productDomainService, ordersBaristaEventPublisher, ordersKitchenEventPublisher)
	counterServiceServer := router.NewGRPCCounterServer(grpcServer, cfg, useCase)
	baristaOrderUpdatedEventHandler := handlers.NewBaristaOrderUpdatedEventHandler(orderRepo)
	kitchenOrderUpdatedEventHandler := handlers.NewKitchenOrderUpdatedEventHandler(orderRepo)
	app := New(cfg, dbEngine, connection, ordersBaristaEventPublisher, ordersKitchenEventPublisher, eventConsumer, productDomainService, useCase, counterServiceServer, baristaOrderUpdatedEventHandler, kitchenOrderUpdatedEventHandler)
	return app, nil
}

// wire.go:

func baristaEventPublisher(amqpConn *amqp091.Connection) orders.BaristaEventPublisher {
	pub, _ := publisher.NewPublisher(amqpConn)
	return (orders.BaristaEventPublisher)(pub)
}

func kitchenEventPublisher(amqpConn *amqp091.Connection) orders.KitchenEventPublisher {
	pub, _ := publisher.NewPublisher(amqpConn)
	return (orders.KitchenEventPublisher)(pub)
}
