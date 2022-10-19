package app

import (
	"context"
	"net"

	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	gen "github.com/thangchung/go-coffeeshop/proto"
	"google.golang.org/grpc"
)

type App struct {
	logger  *mylogger.Logger
	cfg     *config.Config
	network string
	address string
}

type ProductServiceServerImpl struct {
	gen.UnimplementedProductServiceServer
	logger *mylogger.Logger
}

func (g *ProductServiceServerImpl) GetItemTypes(ctx context.Context, request *gen.GetItemTypesRequest) (*gen.GetItemTypesResponse, error) {
	g.logger.Info("GET: GetItemTypes")

	itemTypes := []gen.ItemType{
		{
			Name: "CAPPUCCINO",
			Type: 0,
		},
		{
			Name: "COFFEE_BLACK",
			Type: 1,
		},
		{
			Name: "COFFEE_WITH_ROOM",
			Type: 2,
		},
		{
			Name: "ESPRESSO",
			Type: 3,
		},
		{
			Name: "ESPRESSO_DOUBLE",
			Type: 4,
		},
		{
			Name: "LATTE",
			Type: 5,
		},
		{
			Name: "CAKEPOP",
			Type: 6,
		},
		{
			Name: "CROISSANT",
			Type: 7,
		},
		{
			Name: "MUFFIN",
			Type: 8,
		},
		{
			Name: "CROISSANT_CHOCOLATE",
			Type: 9,
		},
	}

	res := gen.GetItemTypesResponse{}

	for _, v := range itemTypes {
		res.ItemTypes = append(res.ItemTypes, &gen.ItemType{
			Name: v.Name,
			Type: v.Type,
		})
	}

	return &res, nil
}

func New(log *mylogger.Logger, cfg *config.Config) *App {
	return &App{
		logger:  log,
		cfg:     cfg,
		network: "tcp",
		address: "0.0.0.0:5001",
	}
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Init %s %s\n", a.cfg.Name, a.cfg.Version)

	// Repository
	// ...

	// Use case
	// ...

	// gRPC Server
	l, err := net.Listen(a.network, a.address)
	if err != nil {
		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			a.logger.Error("Failed to close %s %s: %v", a.network, a.address, err)
		}
	}()

	s := grpc.NewServer()
	gen.RegisterProductServiceServer(s, &ProductServiceServerImpl{logger: a.logger})

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	a.logger.Info("Start server at " + a.address + " ...")

	return s.Serve(l)
}
