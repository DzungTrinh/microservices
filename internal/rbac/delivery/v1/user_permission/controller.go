package user_permission

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/user_permission"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
	"time"
)

type UserPermissionController struct {
	uc user_permission.UserPermissionUseCase
}

func NewUserPermissionController(uc user_permission.UserPermissionUseCase) *UserPermissionController {
	return &UserPermissionController{uc: uc}
}

func (c *UserPermissionController) AssignPermissionsToUser(ctx context.Context, req *rbacv1.AssignPermissionsToUserRequest) (*rbacv1.AssignPermissionsToUserResponse, error) {
	// Validate UUIDs
	_, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.GetInstance().Errorf("Invalid User ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}
	for _, permID := range req.PermissionIds {
		_, err := uuid.Parse(permID)
		if err != nil {
			logger.GetInstance().Errorf("Invalid Permission ID format: %v", permID)
			return nil, status.Errorf(codes.InvalidArgument, "Invalid Permission ID format: %s", permID)
		}
	}
	if req.GranterId != "" {
		_, err = uuid.Parse(req.GranterId)
		if err != nil {
			logger.GetInstance().Errorf("Invalid Granter ID format: %v", err)
			return nil, status.Errorf(codes.InvalidArgument, "Invalid Granter ID format")
		}
	}

	var expiresAt time.Time
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			logger.GetInstance().Errorf("Invalid expires_at format: %v", err)
			return nil, status.Errorf(codes.InvalidArgument, "Invalid expires_at format")
		}
		expiresAt = t
	} else {
		expiresAt = constants.FallbackFutureTime // no expiration
	}
	dtos := make([]domain.UserPermission, len(req.PermissionIds))
	for i, permID := range req.PermissionIds {
		dtos[i] = domain.UserPermission{
			UserID:       req.UserId,
			PermissionID: permID,
			GranterID:    req.GranterId,
			ExpiresAt:    expiresAt,
		}
	}
	err = c.uc.AssignPermissionsToUser(ctx, dtos)
	if err != nil {
		logger.GetInstance().Errorf("AssignPermissionsToUser failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &rbacv1.AssignPermissionsToUserResponse{Success: true}, nil
}

func (c *UserPermissionController) ListPermissionsForUser(ctx context.Context, req *rbacv1.ListPermissionsForUserRequest) (*rbacv1.ListPermissionsForUserResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.GetInstance().Errorf("Invalid User ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}

	resp, err := c.uc.ListPermissionsForUser(ctx, req.UserId)
	if err != nil {
		logger.GetInstance().Errorf("ListPermissionsForUser failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	pbPerms := make([]*rbacv1.Permission, len(resp))
	for i, p := range resp {
		pbPerms[i] = &rbacv1.Permission{Id: p.ID, Name: p.Name, CreatedAt: p.CreatedAt.Format(time.RFC3339), DeletedAt: ""}
	}
	return &rbacv1.ListPermissionsForUserResponse{Permissions: pbPerms, Success: true}, nil
}

func (c *UserPermissionController) RemovePermissionFromUser(ctx context.Context, req *rbacv1.RemovePermissionFromUserRequest) (*rbacv1.RemovePermissionFromUserResponse, error) {
	_, err := c.uc.RemovePermissionFromUser(ctx, domain.UserPermission{
		UserID:       req.UserId,
		PermissionID: req.PermissionId,
	})
	if err != nil {
		logger.GetInstance().Errorf("RemovePermissionFromUser failed: %v", err)
		return &rbacv1.RemovePermissionFromUserResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &rbacv1.RemovePermissionFromUserResponse{
		Success: true,
	}, nil
}
