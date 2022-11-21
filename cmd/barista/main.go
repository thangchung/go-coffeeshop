package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/thangchung/go-coffeeshop/cmd/barista/config"
	"github.com/thangchung/go-coffeeshop/internal/barista/app"
	mylog "github.com/thangchung/go-coffeeshop/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		glog.Fatal(err)
	}

	logger := mylog.New(cfg.Level)

	a := app.New(logger, cfg)
	if err = a.Run(); err != nil {
		glog.Fatal(err)
		os.Exit(1)
	}
}
