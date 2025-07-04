package role_permission

import (
	"context"
	"errors"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/pkg/logger"
)

type rolePermissionService struct {
	repo repo.RolePermissionRepository
}

func NewRolePermissionService(repo repo.RolePermissionRepository) RolePermissionUseCase {
	return &rolePermissionService{repo: repo}
}

func (s *rolePermissionService) AssignPermissionsToRole(ctx context.Context, rolePerms []domain.RolePermission) error {
	for _, rp := range rolePerms {
		err := s.repo.AssignPermissionsToRole(ctx, rp)
		if err != nil {
			logger.GetInstance().Errorf("Failed to assign permission %s to role %s: %v", rp.PermissionID, rp.RoleID, err)
			return err
		}
	}
	logger.GetInstance().Infof("Assigned %d permissions to role %s", len(rolePerms), rolePerms[0].RoleID)
	return nil
}

func (s *rolePermissionService) ListPermissionsForRole(ctx context.Context, roleID string) ([]domain.Permission, error) {
	perms, err := s.repo.ListPermissionsForRole(ctx, roleID)
	if err != nil {
		logger.GetInstance().Errorf("Failed to list permissions for role %s: %v", roleID, err)
		return nil, err
	}
	logger.GetInstance().Infof("Listed %d permissions for role %s", len(perms), roleID)
	return perms, nil
}

func (s *rolePermissionService) RemovePermissionFromRole(ctx context.Context, rolePerm domain.RolePermission) (*domain.RolePermission, error) {
	if rolePerm.RoleID == "" || rolePerm.PermissionID == "" {
		return nil, errors.New("role_id and permission_id are required")
	}

	err := s.repo.SoftDeleteRolePermission(ctx, rolePerm)
	if err != nil {
		logger.GetInstance().Errorf("RemovePermissionFromRole failed: %v", err)
		return nil, err
	}

	return &rolePerm, nil
}
