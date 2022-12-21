package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/app"
	"github.com/thangchung/go-coffeeshop/pkg/logger"
	"golang.org/x/exp/slog"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed get config", err)
	}

	// set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	// integrate Logrus with the slog logger
	slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))

	a := app.New(cfg)
	if err = a.Run(); err != nil {
		slog.Error("failed app run", err)
		os.Exit(1)
	}
}
