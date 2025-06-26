package role

import (
	"context"
	"github.com/google/uuid"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/pkg/logger"
)

type roleService struct {
	repo repo.RoleRepository
}

func NewRoleService(repo repo.RoleRepository) RoleUseCase {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	role.ID = uuid.New().String()
	created, err := s.repo.CreateRole(ctx, *role)
	if err != nil {
		logger.GetInstance().Errorf("Failed to create role %s: %v", role.Name, err)
		return nil, err
	}
	logger.GetInstance().Infof("Role created: %s (ID: %s)", created.Name, created.ID)
	return &created, nil
}

func (s *roleService) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	role, err := s.repo.GetRoleByName(ctx, name)
	if err != nil {
		logger.GetInstance().Errorf("Failed to get role %s: %v", name, err)
		return nil, err
	}
	return &role, nil
}

func (s *roleService) ListRoles(ctx context.Context) ([]domain.Role, error) {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		logger.GetInstance().Errorf("Failed to list roles: %v", err)
		return nil, err
	}
	logger.GetInstance().Infof("Listed %d roles", len(roles))
	return roles, nil
}

func (s *roleService) UpdateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	updated, err := s.repo.UpdateRole(ctx, *role)
	if err != nil {
		logger.GetInstance().Errorf("Failed to update role %s: %v", role.ID, err)
		return nil, err
	}
	logger.GetInstance().Infof("Role updated: %s (ID: %s)", updated.Name, updated.ID)
	return &updated, nil
}

func (s *roleService) DeleteRole(ctx context.Context, id string) error {
	err := s.repo.DeleteRole(ctx, id)
	if err != nil {
		logger.GetInstance().Errorf("Failed to delete role %s: %v", id, err)
		return err
	}
	logger.GetInstance().Infof("Role deleted: %s", id)
	return nil
}
