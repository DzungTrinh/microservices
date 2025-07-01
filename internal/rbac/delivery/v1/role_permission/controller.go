package role_permission

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/role_permission"
	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
	"time"
)

type RolePermissionController struct {
	uc role_permission.RolePermissionUseCase
}

func NewRolePermissionController(uc role_permission.RolePermissionUseCase) *RolePermissionController {
	return &RolePermissionController{uc: uc}
}

func (c *RolePermissionController) AssignPermissionsToRole(ctx context.Context, req *rbacv1.AssignPermissionsToRoleRequest) (*rbacv1.AssignPermissionsToRoleResponse, error) {
	// Validate UUIDs
	_, err := uuid.Parse(req.RoleId)
	if err != nil {
		logger.GetInstance().Errorf("Invalid Role ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Role ID format")
	}
	for _, permID := range req.PermissionIds {
		_, err := uuid.Parse(permID)
		if err != nil {
			logger.GetInstance().Errorf("Invalid Permission ID format: %v", permID)
			return nil, status.Errorf(codes.InvalidArgument, "Invalid Permission ID format: %s", permID)
		}
	}

	dtos := make([]domain.RolePermission, len(req.PermissionIds))
	for i, permID := range req.PermissionIds {
		dtos[i] = domain.RolePermission{RoleID: req.RoleId, PermissionID: permID}
	}
	err = c.uc.AssignPermissionsToRole(ctx, dtos)
	if err != nil {
		logger.GetInstance().Errorf("AssignPermissionsToRole failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &rbacv1.AssignPermissionsToRoleResponse{Success: true}, nil
}

func (c *RolePermissionController) ListPermissionsForRole(ctx context.Context, req *rbacv1.ListPermissionsForRoleRequest) (*rbacv1.ListPermissionsForRoleResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.RoleId)
	if err != nil {
		logger.GetInstance().Errorf("Invalid Role ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Role ID format")
	}

	resp, err := c.uc.ListPermissionsForRole(ctx, req.RoleId)
	if err != nil {
		logger.GetInstance().Errorf("ListPermissionsForRole failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	pbPerms := make([]*rbacv1.Permission, len(resp))
	for i, p := range resp {
		pbPerms[i] = &rbacv1.Permission{Id: p.ID, Name: p.Name, CreatedAt: p.CreatedAt.Format(time.RFC3339), DeletedAt: ""}
	}
	return &rbacv1.ListPermissionsForRoleResponse{Permissions: pbPerms, Success: true}, nil
}

func (c *RolePermissionController) RemovePermissionFromRole(ctx context.Context, req *rbacv1.RemovePermissionFromRoleRequest) (*rbacv1.RemovePermissionFromRoleResponse, error) {
	_, err := c.uc.RemovePermissionFromRole(ctx, domain.RolePermission{
		RoleID:       req.RoleId,
		PermissionID: req.PermissionId,
	})
	if err != nil {
		logger.GetInstance().Errorf("RemovePermissionFromRole failed: %v", err)
		return &rbacv1.RemovePermissionFromRoleResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &rbacv1.RemovePermissionFromRoleResponse{
		Success: true,
	}, nil
}
