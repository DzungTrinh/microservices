package user_role

import (
	"context"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/pkg/logger"
)

type userRoleService struct {
	repo   repo.UserRoleRepository
	roleUC role.RoleUseCase
}

func NewUserRoleService(repo repo.UserRoleRepository, roleUC role.RoleUseCase) UserRoleUseCase {
	return &userRoleService{
		repo:   repo,
		roleUC: roleUC,
	}
}

func (s *userRoleService) AssignRolesToUser(ctx context.Context, userRoles []domain.UserRole) error {

	for i, userRole := range userRoles {
		if s.roleUC == nil {
			logger.GetInstance().Error("roleUC is nil in AssignRolesToUser")
			panic("roleUC is nil in AssignRolesToUser")
		}

		// Fetch role ID by name
		role, err := s.roleUC.GetRoleByName(ctx, userRole.RoleName)
		if err != nil {
			logger.GetInstance().Errorf("Failed to get role %s for user %s: %v", userRole.RoleID, userRole.UserID, err)
			return err
		}
		// Update userRole with actual role ID
		userRoles[i].RoleID = role.ID
	}

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
