package main

import (
	"context"
	"os"

	"github.com/golang/glog"
	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	"github.com/thangchung/go-coffeeshop/internal/product/app"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		glog.Fatal(err)
	}

	mylog := mylogger.New(cfg.Level)

	a := app.New(mylog, cfg)
	if err = a.Run(context.Background()); err != nil {
		glog.Fatal(err)
		os.Exit(1)
	}
}
