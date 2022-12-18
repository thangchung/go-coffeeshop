package app

import (
	"context"
	"fmt"
	"net"

	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	productGrpc "github.com/thangchung/go-coffeeshop/internal/product/infras/grpc"
	productRepo "github.com/thangchung/go-coffeeshop/internal/product/infras/repo"
	productUC "github.com/thangchung/go-coffeeshop/internal/product/usecases/products"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

type App struct {
	cfg     *config.Config
	network string
	address string
}

func New(cfg *config.Config) *App {
	return &App{
		cfg:     cfg,
		network: "tcp",
		address: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
	}
}

func (a *App) Run() error {
	slog.Info("init", "name", a.cfg.Name, "version", a.cfg.Version)

	ctx, _ := context.WithCancel(context.Background())

	// Repository
	repo := productRepo.NewOrderRepo()

	// UC
	uc := productUC.NewService(repo)

	// gRPC Server
	l, err := net.Listen(a.network, a.address)
	if err != nil {
		slog.Error("failed app-Run-net.Listener", err)

		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			slog.Error("failed to close", err, "network", a.network, "address", a.address)
		}
	}()

	s := grpc.NewServer()
	productGrpc.NewProductGRPCServer(s, uc)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	slog.Info("start server", "address", a.address)

	return s.Serve(l)
}
