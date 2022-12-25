package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	configs "github.com/thangchung/go-coffeeshop/pkg/config"
)

type (
	Config struct {
		configs.App   `yaml:"app"`
		configs.HTTP  `yaml:"http"`
		configs.Log   `yaml:"logger"`
		PG            `yaml:"postgres"`
		RabbitMQ      `yaml:"rabbitmq"`
		ProductClient `yaml:"product_client"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		DsnURL  string `env-required:"true" yaml:"dsn_url" env:"PG_DSN_URL"`
	}

	RabbitMQ struct {
		URL string `env-required:"true" yaml:"url" env:"RABBITMQ_URL"`
	}

	ProductClient struct {
		URL string `env-required:"true" yaml:"url" env:"PRODUCT_CLIENT_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println("config path: " + dir)

	err = cleanenv.ReadConfig(dir+"/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
