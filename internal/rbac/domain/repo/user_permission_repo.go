package repo

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type UserPermissionRepository interface {
	AssignPermissionsToUser(ctx context.Context, userPerm domain.UserPermission) error
	ListPermissionsForUser(ctx context.Context, userID string) ([]domain.Permission, error)
}
