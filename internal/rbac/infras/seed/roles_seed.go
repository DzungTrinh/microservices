package seed

import (
	"context"
	"github.com/google/uuid"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
)

func SeedRoles(ctx context.Context, roleUC role.RoleUseCase) error {
	roles := []struct {
		name string
	}{
		{name: constants.RoleUser},
		{name: constants.RoleAdmin},
	}

	for _, r := range roles {
		role := domain.Role{
			ID:      uuid.New().String(),
			Name:    r.name,
			BuiltIn: true,
		}

		_, err := roleUC.CreateRole(ctx, &role)
		if err != nil {
			if err.Error() == "role already exists" { // Adjust based on your RoleUseCase error
				logger.GetInstance().Printf("Role %s already exists", r.name)
				continue
			}
			logger.GetInstance().Errorf("Failed to seed role %s: %v", r.name, err)
			return err
		}
		logger.GetInstance().Printf("Role %s created successfully", r.name)
	}

	return nil
}
