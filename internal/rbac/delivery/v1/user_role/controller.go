package user_role

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/user_role"
	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
)

type UserRoleController struct {
	uc user_role.UserRoleUseCase
}

func NewUserRoleController(uc user_role.UserRoleUseCase) *UserRoleController {
	return &UserRoleController{uc: uc}
}

func (c *UserRoleController) AssignRolesToUser(ctx context.Context, req *rbacv1.AssignRolesToUserRequest) (*rbacv1.AssignRolesToUserResponse, error) {
	// Validate UUIDs
	_, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.GetInstance().Errorf("Invalid User ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}
	for _, roleID := range req.RoleIds {
		_, err := uuid.Parse(roleID)
		if err != nil {
			logger.GetInstance().Errorf("Invalid Role ID format: %v", roleID)
			return nil, status.Errorf(codes.InvalidArgument, "Invalid Role ID format: %s", roleID)
		}
	}

	dtos := make([]domain.UserRole, len(req.RoleIds))
	for i, roleID := range req.RoleIds {
		dtos[i] = domain.UserRole{UserID: req.UserId, RoleID: roleID}
	}
	err = c.uc.AssignRolesToUser(ctx, dtos)
	if err != nil {
		logger.GetInstance().Errorf("AssignRolesToUser failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &rbacv1.AssignRolesToUserResponse{Success: true}, nil
}

func (c *UserRoleController) ListRolesForUser(ctx context.Context, req *rbacv1.ListRolesForUserRequest) (*rbacv1.ListRolesForUserResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.GetInstance().Errorf("Invalid User ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User ID format")
	}

	resp, err := c.uc.ListRolesForUser(ctx, req.UserId)
	if err != nil {
		logger.GetInstance().Errorf("ListRolesForUser failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	pbRoles := make([]*rbacv1.Role, len(resp))
	for i, r := range resp {
		pbRoles[i] = &rbacv1.Role{Id: r.ID, Name: r.Name, BuiltIn: r.BuiltIn}
	}
	return &rbacv1.ListRolesForUserResponse{Roles: pbRoles, Success: true}, nil
}
