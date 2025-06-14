package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseDSN   string `env:"DATABASE_DSN" env-required:"true"`
	JWTSecret     string `env:"JWT_SECRET" env-required:"true"`
	RefreshSecret string `env:"REFRESH_SECRET" env-required:"true"`
	Port          string `env:"PORT" env-required:"true"`
	AdminEmail    string `env:"ADMIN_EMAIL" env-required:"true"`
	AdminPassword string `env:"ADMIN_PASSWORD" env-required:"true"`
}

func Load() Config {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}
	return cfg
}
