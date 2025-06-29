package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseDSN string `env:"DATABASE_DSN" env-required:"true"`
	JWTSecret   string `env:"JWT_SECRET" env-required:"true"`
	Port        string `env:"PORT" env-required:"true"`
}

func Load() Config {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		logger.GetInstance().Fatalf("Failed to load .env: %v", err)
	}
	return cfg
}
