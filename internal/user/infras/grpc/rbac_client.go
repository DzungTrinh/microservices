package grpc

import (
	"context"
	"microservices/user-management/internal/user/usecases/grpc"

	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
)

type rbacService struct {
	rbacClient rbacv1.RBACServiceClient
}

func NewRBACService(rbacClient rbacv1.RBACServiceClient) grpc.RBACService {
	return &rbacService{rbacClient: rbacClient}
}

func (s *rbacService) ListRolesForUser(ctx context.Context, userID string) ([]string, error) {
	resp, err := s.rbacClient.ListRolesForUser(ctx, &rbacv1.ListRolesForUserRequest{UserId: userID})
	if err != nil {
		logger.GetInstance().Errorf("Failed to fetch roles for user %s: %v", userID, err)
		return nil, err
	}
	roles := make([]string, len(resp.Roles))
	for i, role := range resp.Roles {
		roles[i] = role.Name
	}
	return roles, nil
}

func (s *rbacService) ListPermissionsForUser(ctx context.Context, userID string) ([]string, error) {
	resp, err := s.rbacClient.ListPermissionsForUser(ctx, &rbacv1.ListPermissionsForUserRequest{UserId: userID})
	if err != nil {
		logger.GetInstance().Errorf("Failed to fetch permissions for user %s: %v", userID, err)
		return nil, err
	}
	permissions := make([]string, len(resp.Permissions))
	for i, perm := range resp.Permissions {
		permissions[i] = perm.Name
	}
	return permissions, nil
}
