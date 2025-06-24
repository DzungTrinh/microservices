package config

import (
	"go.uber.org/zap"
	logger "microservices/user-management/pkg/logger/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseDSN   string        `env:"DATABASE_DSN" env-required:"true"`
	JWTSecret     string        `env:"JWT_SECRET" env-required:"true"`
	RefreshSecret string        `env:"REFRESH_SECRET" env-required:"true"`
	Port          string        `env:"PORT" env-required:"true"`
	AdminEmail    string        `env:"ADMIN_EMAIL" env-required:"true"`
	AdminPassword string        `env:"ADMIN_PASSWORD" env-required:"true"`
	GRPCPort      string        `env:"GRPC_PORT" env-required:"true"`
	Logger        logger.Config `env:"LOGGER"`
}

func Load() error {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		zap.S().Error(err)
	}
	instance = &cfg
	return nil
}
