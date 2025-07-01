package user_permission

import (
	"context"
	"errors"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/pkg/logger"
)

type userPermissionService struct {
	repo repo.UserPermissionRepository
}

func NewUserPermissionService(repo repo.UserPermissionRepository) UserPermissionUseCase {
	return &userPermissionService{repo: repo}
}

func (s *userPermissionService) AssignPermissionsToUser(ctx context.Context, userPerms []domain.UserPermission) error {
	for _, up := range userPerms {
		err := s.repo.AssignPermissionsToUser(ctx, up)
		if err != nil {
			logger.GetInstance().Errorf("Failed to assign permission %s to user %s: %v", up.PermissionID, up.UserID, err)
			return err
		}
	}
	logger.GetInstance().Infof("Assigned %d permissions to user %s", len(userPerms), userPerms[0].UserID)
	return nil
}

func (s *userPermissionService) ListPermissionsForUser(ctx context.Context, userID string) ([]domain.Permission, error) {
	perms, err := s.repo.ListPermissionsForUser(ctx, userID)
	if err != nil {
		logger.GetInstance().Errorf("Failed to list permissions for user %s: %v", userID, err)
		return nil, err
	}
	logger.GetInstance().Infof("Listed %d permissions for user %s", len(perms), userID)
	return perms, nil
}

func (s *userPermissionService) RemovePermissionFromUser(ctx context.Context, userPerm domain.UserPermission) (*domain.UserPermission, error) {
	if userPerm.UserID == "" || userPerm.PermissionID == "" {
		return nil, errors.New("user_id and permission_id are required")
	}

	err := s.repo.SoftDeleteUserPermission(ctx, userPerm)
	if err != nil {
		logger.GetInstance().Errorf("RemovePermissionFromUser failed: %v", err)
		return nil, err
	}

	return &userPerm, nil
}
