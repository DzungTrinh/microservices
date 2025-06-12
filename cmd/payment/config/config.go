package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseDSN string `env:"DATABASE_DSN" env-required:"true"`
	JWTSecret   string `env:"JWT_SECRET" env-required:"true"`
	Port        string `env:"PORT" env-default:"8080"`
}

func Load() Config {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}
	return cfg
}
