package repo

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type PermissionRepository interface {
	CreatePermission(ctx context.Context, perm domain.Permission) error
	DeletePermission(ctx context.Context, id string) error
	ListPermissions(ctx context.Context) ([]domain.Permission, error)
}
