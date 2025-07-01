package user_role

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type UserRoleUseCase interface {
	AssignRolesToUser(ctx context.Context, userRoles []domain.UserRole) error
	ListRolesForUser(ctx context.Context, userID string) ([]domain.Role, error)
	RemoveRoleFromUser(ctx context.Context, userRole domain.UserRole) (*domain.UserRole, error)
}
