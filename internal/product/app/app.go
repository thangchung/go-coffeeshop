package app

import (
	"context"
	"fmt"
	"net"

	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	productRepo "github.com/thangchung/go-coffeeshop/internal/product/features/products/repo"
	productGrpc "github.com/thangchung/go-coffeeshop/internal/product/grpc"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	logger  *mylogger.Logger
	cfg     *config.Config
	network string
	address string
}

func New(log *mylogger.Logger, cfg *config.Config) *App {
	return &App{
		logger:  log,
		cfg:     cfg,
		network: "tcp",
		address: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
	}
}

func (a *App) Run() error {
	a.logger.Info("Init %s %s\n", a.cfg.Name, a.cfg.Version)

	ctx, _ := context.WithCancel(context.Background())

	// Repository
	repo := productRepo.NewOrderRepo()

	// gRPC Server
	l, err := net.Listen(a.network, a.address)
	if err != nil {
		a.logger.Fatal("app-Run-net.Listener: %s", err.Error())

		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			a.logger.Error("Failed to close %s %s: %v", a.network, a.address, err)
		}
	}()

	s := grpc.NewServer()
	productGrpc.NewProductServiceServerGrpc(s, a.logger, repo)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	a.logger.Info("Start server at " + a.address + " ...")

	return s.Serve(l)
}
