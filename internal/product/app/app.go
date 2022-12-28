package app

import (
	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	productUC "github.com/thangchung/go-coffeeshop/internal/product/usecases/products"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type App struct {
	Cfg               *config.Config
	UC                productUC.UseCase
	ProductGRPCServer gen.ProductServiceServer
}

func New(
	cfg *config.Config,
	uc productUC.UseCase,
	productGRPCServer gen.ProductServiceServer,
) *App {
	return &App{
		Cfg:               cfg,
		UC:                uc,
		ProductGRPCServer: productGRPCServer,
	}
}
