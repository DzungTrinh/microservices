package role

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
)

type RoleUseCase interface {
	CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error)
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	ListRoles(ctx context.Context) ([]domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) (*domain.Role, error)
	DeleteRole(ctx context.Context, id string) error
}
