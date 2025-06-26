package seed

import (
	"context"
	"microservices/user-management/internal/user/usecases/user"
	"microservices/user-management/pkg/logger"
)

func SeedAdmin(ctx context.Context, adminUC user.UserUseCase, adminEmail, adminPassword string) error {
	_, err := adminUC.CreateAdmin(ctx, adminEmail, "admin", adminPassword)
	if err != nil {
		if err.Error() == "admin user already exists" {
			logger.GetInstance().Printf("Admin user %s already exists", adminEmail)
			return nil
		}
		logger.GetInstance().Errorf("Failed to seed admin: %v", err)
		return err
	}

	logger.GetInstance().Printf("Admin user %s created successfully", adminEmail)
	return nil
}
