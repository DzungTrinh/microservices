package main

import (
	"log"
	"microservices/user-management/cmd/user/config"
	"microservices/user-management/internal/user/app"
)

func main() {
	cfg := config.Load()

	application := app.NewApp(cfg)
	defer application.Close()

	if err := application.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
