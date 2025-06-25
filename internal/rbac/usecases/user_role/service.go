package user_role

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/pkg/logger"
)

type userRoleService struct {
	repo repo.UserRoleRepository
}

func NewUserRoleService(repo repo.UserRoleRepository) UserRoleUseCase {
	return &userRoleService{repo: repo}
}

func (s *userRoleService) AssignRolesToUser(ctx context.Context, userRoles []domain.UserRole) error {
	for _, ur := range userRoles {
		err := s.repo.AssignRolesToUser(ctx, ur)
		if err != nil {
			logger.GetInstance().Errorf("Failed to assign role %s to user %s: %v", ur.RoleID, ur.UserID, err)
			return err
		}
	}
	logger.GetInstance().Infof("Assigned %d roles to user %s", len(userRoles), userRoles[0].UserID)
	return nil
}

func (s *userRoleService) ListRolesForUser(ctx context.Context, userID string) ([]domain.Role, error) {
	roles, err := s.repo.ListRolesForUser(ctx, userID)
	if err != nil {
		logger.GetInstance().Errorf("Failed to list roles for user %s: %v", userID, err)
		return nil, err
	}
	logger.GetInstance().Infof("Listed %d roles for user %s", len(roles), userID)
	return roles, nil
}
