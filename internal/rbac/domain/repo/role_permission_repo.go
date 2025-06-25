package repo

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type RolePermissionRepository interface {
	AssignPermissionsToRole(ctx context.Context, rolePerm domain.RolePermission) error
	ListPermissionsForRole(ctx context.Context, roleID string) ([]domain.Permission, error)
}
