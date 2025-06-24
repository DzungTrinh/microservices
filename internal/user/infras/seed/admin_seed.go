package seed

import (
	"context"
	"microservices/user-management/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/infras/mysql"
)

// SeedAdmin creates an admin user if it doesn't exist.
func SeedAdmin(ctx context.Context, queries *mysql.Queries, adminEmail, adminPassword string) error {
	// Check if admin exists
	_, err := queries.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		logger.GetInstance().Printf("Admin user %s already exists", adminEmail)
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin user
	userID := uuid.New().String()
	_, err = queries.CreateUser(ctx, mysql.CreateUserParams{
		ID:       userID,
		Username: "admin",
		Email:    adminEmail,
		Password: string(hashedPassword),
	})
	if err != nil {
		return err
	}

	// Assign admin role
	roleID, err := queries.GetRoleIDByName(ctx, string(domain.RoleAdmin))
	if err != nil {
		return err
	}

	err = queries.CreateUserRole(ctx, mysql.CreateUserRoleParams{
		UserID: userID,
		RoleID: roleID,
	})
	if err != nil {
		return err
	}

	logger.GetInstance().Printf("Admin user %s created successfully", adminEmail)
	return nil
}
