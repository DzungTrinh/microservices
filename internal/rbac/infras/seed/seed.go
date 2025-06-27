package seed

import (
	"context"
	"github.com/google/uuid"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/permission"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/internal/rbac/usecases/role_permission"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
)

func SeedRoles(ctx context.Context, roleUC role.RoleUseCase, rpUC role_permission.RolePermissionUseCase, permUC permission.PermissionUseCase) error {
	// Define roles and their permissions
	roles := []struct {
		name        string
		permissions []string
	}{
		{
			name:        constants.RoleAdmin,
			permissions: []string{"read_profile", "write_profile", "manage_users", "manage_roles", "manage_permissions"},
		},
		{
			name:        constants.RoleUser,
			permissions: []string{"read_profile", "write_profile"},
		},
	}

	// Map to store permission IDs for assignment
	permissionIDs := make(map[string]string)

	// Seed permissions
	permissions := []string{
		"read_profile",
		"write_profile",
		"manage_users",
		"manage_roles",
		"manage_permissions",
	}

	for _, permName := range permissions {
		perm := domain.Permission{
			Name: permName,
		}
		permID, err := permUC.CreatePermission(ctx, &perm)
		if err != nil {
			logger.GetInstance().Errorf("Failed to seed permission %s: %v", permName, err)
			continue
		}
		permissionIDs[permName] = permID
		logger.GetInstance().Printf("Permission %s created successfully", permName)
	}

	// Seed roles and assign permissions
	for _, r := range roles {
		role := domain.Role{
			ID:      uuid.New().String(),
			Name:    r.name,
			BuiltIn: true,
		}

		createdRole, err := roleUC.CreateRole(ctx, &role)
		if err != nil {
			logger.GetInstance().Errorf("Failed to seed role %s: %v", r.name, err)
			continue
		}
		role.ID = createdRole.ID
		logger.GetInstance().Printf("Role %s created successfully", r.name)

		// Assign permissions to role
		var rolePerms []domain.RolePermission
		for _, permName := range r.permissions {
			if permID, exists := permissionIDs[permName]; exists {
				rolePerms = append(rolePerms, domain.RolePermission{
					RoleID:       role.ID,
					PermissionID: permID,
				})
			}
		}

		if len(rolePerms) > 0 {
			err = rpUC.AssignPermissionsToRole(ctx, rolePerms)
			if err != nil {
				logger.GetInstance().Errorf("Failed to assign permissions to role %s: %v", r.name, err)
				continue
			}
			logger.GetInstance().Printf("Assigned permissions %v to role %s", r.permissions, r.name)
		}
	}

	return nil
}
