package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	"github.com/thangchung/go-coffeeshop/internal/product/app"
	"github.com/thangchung/go-coffeeshop/pkg/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func main() {
	// set GOMAXPROCS
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("failed set max procs", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed get config", err)
	}

	slog.Info("⚡ init app", "name", cfg.Name, "version", cfg.Version)

	// set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	// integrate Logrus with the slog logger
	slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))

	server := grpc.NewServer()

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	_, err = app.InitApp(cfg, server)
	if err != nil {
		slog.Error("failed init app", err)
		cancel()
	}

	// gRPC Server.
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("failed to listen to address", err, "network", network, "address", address)
		cancel()
	}

	slog.Info("🌏 start server...", "address", address)

	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close", err1, "network", network, "address", address)
		}
	}()

	err = server.Serve(l)
	if err != nil {
		slog.Error("failed start gRPC server", err, "network", network, "address", address)
		cancel()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		slog.Info("signal.Notify", v)
	case done := <-ctx.Done():
		slog.Info("ctx.Done", done)
	}
}
