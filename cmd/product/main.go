package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	log "github.com/thangchung/go-coffeeshop/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Config error: %s", err)
	}

	logger := log.New(cfg.Log.Level)

	logger.Info("Hello %s %s\n", cfg.Name, cfg.Version)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
