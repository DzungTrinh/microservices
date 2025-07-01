package role_permission

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type RolePermissionUseCase interface {
	AssignPermissionsToRole(ctx context.Context, rolePerms []domain.RolePermission) error
	ListPermissionsForRole(ctx context.Context, roleID string) ([]domain.Permission, error)
	RemovePermissionFromRole(ctx context.Context, rolePerm domain.RolePermission) (*domain.RolePermission, error)
}
