package permission

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/permission"
	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
)

type PermissionController struct {
	uc permission.PermissionUseCase
}

func NewPermissionController(uc permission.PermissionUseCase) *PermissionController {
	return &PermissionController{uc: uc}
}

func (c *PermissionController) CreatePermission(ctx context.Context, req *rbacv1.CreatePermissionRequest) (*rbacv1.CreatePermissionResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.Name)
	if err == nil {
		logger.GetInstance().Errorf("Invalid Permission Name, Name cannot be UUID")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Permission Name, Name cannot be UUID")
	}

	dto := domain.Permission{
		Name: req.Name,
	}
	err = c.uc.CreatePermission(ctx, &dto)
	if err != nil {
		logger.GetInstance().Errorf("CreatePermission failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &rbacv1.CreatePermissionResponse{Success: true}, nil
}

func (c *PermissionController) DeletePermission(ctx context.Context, req *rbacv1.DeletePermissionRequest) (*rbacv1.DeletePermissionResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.Id)
	if err != nil {
		logger.GetInstance().Errorf("Invalid Permission ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Permission ID format")
	}

	err = c.uc.DeletePermission(ctx, req.Id)
	if err != nil {
		logger.GetInstance().Errorf("DeletePermission failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &rbacv1.DeletePermissionResponse{Success: true}, nil
}

func (c *PermissionController) ListPermissions(ctx context.Context, _ *rbacv1.Empty) (*rbacv1.ListPermissionsResponse, error) {
	resp, err := c.uc.ListPermissions(ctx)
	if err != nil {
		logger.GetInstance().Errorf("ListPermissions failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbPerms := make([]*rbacv1.Permission, len(resp))
	for i, p := range resp {
		pbPerms[i] = &rbacv1.Permission{Id: p.ID, Name: p.Name}
	}
	return &rbacv1.ListPermissionsResponse{Permissions: pbPerms, Success: true}, nil
}
