package permission

import (
	"context"
	"github.com/google/uuid"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/pkg/logger"
)

type permissionService struct {
	repo repo.PermissionRepository
}

func NewPermissionService(repo repo.PermissionRepository) PermissionUseCase {
	return &permissionService{repo: repo}
}

func (s *permissionService) CreatePermission(ctx context.Context, perm *domain.Permission) (string, error) {
	perm.ID = uuid.New().String()
	err := s.repo.CreatePermission(ctx, *perm)
	if err != nil {
		logger.GetInstance().Errorf("Failed to create permission %s: %v", perm.Name, err)
		return "", err
	}
	logger.GetInstance().Infof("Permission created: %s", perm.ID)
	return perm.ID, nil
}

func (s *permissionService) DeletePermission(ctx context.Context, id string) error {
	err := s.repo.DeletePermission(ctx, id)
	if err != nil {
		logger.GetInstance().Errorf("Failed to delete permission %s: %v", id, err)
		return err
	}
	logger.GetInstance().Infof("Permission deleted: %s", id)
	return nil
}

func (s *permissionService) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	permissions, err := s.repo.ListPermissions(ctx)
	if err != nil {
		logger.GetInstance().Errorf("Failed to list permissions: %v", err)
		return nil, err
	}
	logger.GetInstance().Infof("Listed %d permissions", len(permissions))
	return permissions, nil
}
