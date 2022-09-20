package main

import (
	"fmt"
	"log"

	"github.com/thangchung/go-coffeeshop/cmd/product/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	fmt.Printf("Hello %s %s\n", cfg.Name, cfg.Version)
}
