package repo

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type UserRoleRepository interface {
	AssignRolesToUser(ctx context.Context, userRole domain.UserRole) error
	ListRolesForUser(ctx context.Context, userID string) ([]domain.Role, error)
	SoftDeleteUserRole(ctx context.Context, userRole domain.UserRole) error
}
