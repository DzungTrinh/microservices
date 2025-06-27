package config

import (
	"go.uber.org/zap"
	logger "microservices/user-management/pkg/logger/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port         string        `env:"PORT" env-required:"true"`
	UserHttpPort string        `env:"USER_HTTP_PORT" env-required:"true"`
	UserGrpcPort string        `env:"USER_GRPC_PORT" env-required:"true"`
	RbacHttpPort string        `env:"RBAC_HTTP_PORT" env-required:"true"`
	RbacGrpcPort string        `env:"RBAC_GRPC_PORT" env-required:"true"`
	Logger       logger.Config `env:"LOGGER"`
}

func Load() error {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		zap.S().Error(err)
	}
	instance = &cfg
	return nil
}
