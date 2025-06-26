package config

import (
	"go.uber.org/zap"
	logger "microservices/user-management/pkg/logger/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseDSN   string        `env:"DATABASE_DSN" env-required:"true"`
	Port          string        `env:"PORT" env-required:"true"`
	GRPCPort      string        `env:"GRPC_PORT" env-required:"true"`
	Logger        logger.Config `env:"LOGGER"`
	RabbitmqUrl   string        `env:"RABBITMQ_URL" env-required:"true"`
	RabbitmqQueue string        `env:"RABBITMQ_QUEUE" env-required:"true"`
}

func Load() error {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		zap.S().Error(err)
	}
	instance = &cfg
	return nil
}
