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
		configs.App  `yaml:"app"`
		configs.HTTP `yaml:"http"`
		configs.Log  `yaml:"logger"`
		RabbitMQ     `yaml:"rabbit_mq"`
	}

	RabbitMQ struct {
		URL            string `env-required:"true" yaml:"url" env:"RABBITMQ_URL"`
		Exchange       string `env-required:"true" yaml:"exchange" env:"RABBITMQ_Exchange"`
		Queue          string `env-required:"true" yaml:"queue" env:"RABBITMQ_Queue"`
		RoutingKey     string `env-required:"true" yaml:"routing_key" env:"RABBITMQ_RoutingKey"`
		ConsumerTag    string `env-required:"true" yaml:"consumer_tag" env:"RABBITMQ_ConsumerTag"`
		WorkerPoolSize int    `env-required:"true" yaml:"worker_pool_size" env:"RABBITMQ_WorkerPoolSize"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println(dir)

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
