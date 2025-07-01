package permission

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type PermissionUseCase interface {
	CreatePermission(ctx context.Context, perm *domain.Permission) (string, error)
	DeletePermission(ctx context.Context, id string) error
	ListPermissions(ctx context.Context) ([]domain.Permission, error)
}
