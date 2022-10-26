package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/golang/glog"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/app"
	mylogger "github.com/thangchung/go-coffeeshop/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println(reflect.TypeOf(struct{}{}))

	mylog := mylogger.New(cfg.Level)

	a := app.New(mylog, cfg)
	if err = a.Run(context.Background()); err != nil {
		glog.Fatal(err)
		os.Exit(1)
	}
}
