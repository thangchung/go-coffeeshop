package main

import (
	"context"
	"log"

	"github.com/golang/glog"
	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	"github.com/thangchung/go-coffeeshop/internal/product/app"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	mylog := mylogger.New(cfg.Level)

	ctx := context.Background()
	if err = app.Run(ctx, cfg, mylog); err != nil {
		glog.Fatal(err)
	}
}
