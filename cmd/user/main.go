package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"microservices/user-management/internal/user/app"
)

type Config struct {
	DatabaseDSN string `env:"DATABASE_DSN" env-required:"true"`
	JWTSecret   string `env:"JWT_SECRET" env-required:"true"`
	Port        string `env:"PORT" env-default:"8080"`
}

func main() {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	app := app.NewApp(app.Config{
		DatabaseDSN: cfg.DatabaseDSN,
		JWTSecret:   cfg.JWTSecret,
	})
	defer app.Close()
	if err := app.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
