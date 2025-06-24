package main

import (
	"microservices/user-management/cmd/user/config"
	"microservices/user-management/internal/user/app"
	"microservices/user-management/pkg/logger"
)

func main() {
	cfg := config.GetInstance()

	application := app.NewApp(*cfg)
	defer application.Close()

	if err := application.Run(":" + cfg.Port); err != nil {
		logger.GetInstance().Fatalf("Failed to run server: %v", err)
	}
}
