package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	log "github.com/thangchung/go-coffeeshop/pkg/logger"
)

func Run(cfg *config.Config) {
	logger := log.New(cfg.Log.Level)
	logger.Info("Init %s %s\n", cfg.Name, cfg.Version)

	// Repository
	// ...

	// Use case
	// ...

	// HTTP Server
	// ...

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
