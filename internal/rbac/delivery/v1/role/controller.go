package role

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
	"time"
)

type RoleController struct {
	uc role.RoleUseCase
}

func NewRoleController(uc role.RoleUseCase) *RoleController {
	return &RoleController{uc: uc}
}

func (c *RoleController) CreateRole(ctx context.Context, req *rbacv1.CreateRoleRequest) (*rbacv1.CreateRoleResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.Name)
	if err == nil {
		logger.GetInstance().Errorf("Invalid Role Name, Name cannot be UUID")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Role Name, Name cannot be UUID")
	}

	dto := domain.Role{
		Name:    req.Name,
		BuiltIn: req.BuiltIn,
	}
	resp, err := c.uc.CreateRole(ctx, &dto)
	if err != nil {
		logger.GetInstance().Errorf("CreateRole failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &rbacv1.CreateRoleResponse{
		RoleId:  resp.ID,
		Name:    resp.Name,
		Success: true,
	}, nil
}

func (c *RoleController) GetRoleByName(ctx context.Context, req *rbacv1.GetRoleByNameRequest) (*rbacv1.GetRoleByNameResponse, error) {
	resp, err := c.uc.GetRoleByName(ctx, req.Name)
	if err != nil {
		logger.GetInstance().Errorf("GetRoleByName failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &rbacv1.GetRoleByNameResponse{
		RoleId:    resp.ID,
		Name:      resp.Name,
		BuiltIn:   resp.BuiltIn,
		CreatedAt: resp.CreatedAt.Format(time.RFC3339),
		DeletedAt: "",
		Success:   true,
	}, nil
}

func (c *RoleController) ListRoles(ctx context.Context, _ *rbacv1.Empty) (*rbacv1.ListRolesResponse, error) {
	resp, err := c.uc.ListRoles(ctx)
	if err != nil {
		logger.GetInstance().Errorf("ListRoles failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	rbacv1Roles := make([]*rbacv1.Role, len(resp))
	for i, r := range resp {
		rbacv1Roles[i] = &rbacv1.Role{Id: r.ID, Name: r.Name, BuiltIn: r.BuiltIn, CreatedAt: r.CreatedAt.Format(time.RFC3339), DeletedAt: ""}
	}
	return &rbacv1.ListRolesResponse{Roles: rbacv1Roles, Success: true}, nil
}

func (c *RoleController) UpdateRole(ctx context.Context, req *rbacv1.UpdateRoleRequest) (*rbacv1.UpdateRoleResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.Id)
	if err != nil {
		logger.GetInstance().Errorf("Invalid Role ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Role ID format")
	}
	_, err = uuid.Parse(req.Name)
	if err == nil {
		logger.GetInstance().Errorf("Invalid Role Name, Name cannot be UUID")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Role Name, Name cannot be UUID")
	}

	dto := domain.Role{
		ID:      req.Id,
		Name:    req.Name,
		BuiltIn: req.BuiltIn,
	}
	resp, err := c.uc.UpdateRole(ctx, &dto)
	if err != nil {
		logger.GetInstance().Errorf("UpdateRole failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &rbacv1.UpdateRoleResponse{
		RoleId:  resp.ID,
		Name:    resp.Name,
		Success: true,
	}, nil
}

func (c *RoleController) DeleteRole(ctx context.Context, req *rbacv1.DeleteRoleRequest) (*rbacv1.DeleteRoleResponse, error) {
	// Validate UUID
	_, err := uuid.Parse(req.Id)
	if err != nil {
		logger.GetInstance().Errorf("Invalid Role ID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Role ID format")
	}

	err = c.uc.DeleteRole(ctx, req.Id)
	if err != nil {
		logger.GetInstance().Errorf("DeleteRole failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &rbacv1.DeleteRoleResponse{Success: true}, nil
}
