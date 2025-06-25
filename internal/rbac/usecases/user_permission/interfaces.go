package user_permission

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type UserPermissionUseCase interface {
	AssignPermissionsToUser(ctx context.Context, userPerms []domain.UserPermission) error
	ListPermissionsForUser(ctx context.Context, userID string) ([]domain.Permission, error)
}
