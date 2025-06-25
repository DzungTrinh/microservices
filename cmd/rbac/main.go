package main

import (
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/internal/rbac/app"
	"microservices/user-management/pkg/logger"
)

func main() {
	cfg := config.GetInstance()

	application := app.NewApp(*cfg)
	defer func(application *app.App) {
		err := application.Close()
		if err != nil {
			logger.GetInstance().Fatalf("Failed to close server: %v", err)
		}
	}(application)

	if err := application.Run(":" + cfg.Port); err != nil {
		logger.GetInstance().Fatalf("Failed to run server: %v", err)
	}
}
