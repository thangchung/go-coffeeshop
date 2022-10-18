package app

import (
	"context"
	"net"

	"github.com/golang/glog"
	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	gen "github.com/thangchung/go-coffeeshop/proto"
	"google.golang.org/grpc"
)

var (
	logger  *mylogger.Logger
	network = "tcp"
	address = "0.0.0.0:5001"
)

type ProductServiceServerImpl struct {
	gen.UnimplementedProductServiceServer
}

func (g *ProductServiceServerImpl) GetItemTypes(ctx context.Context, request *gen.GetItemTypesRequest) (*gen.GetItemTypesResponse, error) {
	logger.Info("%s", "GET: GetItemTypes")

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

func Run(ctx context.Context, cfg *config.Config, log *mylogger.Logger) error {
	logger = log
	logger.Info("Init %s %s\n", cfg.Name, cfg.Version)

	// Repository
	// ...

	// Use case
	// ...

	// gRPC Server
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			glog.Errorf("Failed to close %s %s: %v", network, address, err)
		}
	}()

	s := grpc.NewServer()
	gen.RegisterProductServiceServer(s, &ProductServiceServerImpl{})

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	return s.Serve(l)
}
