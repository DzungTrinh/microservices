package router

import (
	"context"
	p "microservices/user-management/internal/rbac/delivery/v1/permission"
	r "microservices/user-management/internal/rbac/delivery/v1/role"
	rp "microservices/user-management/internal/rbac/delivery/v1/role_permission"
	up "microservices/user-management/internal/rbac/delivery/v1/user_permission"
	ur "microservices/user-management/internal/rbac/delivery/v1/user_role"
	"microservices/user-management/internal/rbac/usecases/permission"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/internal/rbac/usecases/role_permission"
	"microservices/user-management/internal/rbac/usecases/user_permission"
	"microservices/user-management/internal/rbac/usecases/user_role"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
)

type RBACGrpcServer struct {
	rbacv1.UnimplementedRBACServiceServer
	roleCtrl     *r.RoleController
	permCtrl     *p.PermissionController
	userRoleCtrl *ur.UserRoleController
	userPermCtrl *up.UserPermissionController
	rolePermCtrl *rp.RolePermissionController
}

func NewRBACGrpcServer(
	roleUC role.RoleUseCase,
	permUC permission.PermissionUseCase,
	userRoleUC user_role.UserRoleUseCase,
	userPermUC user_permission.UserPermissionUseCase,
	rolePermUC role_permission.RolePermissionUseCase,
) *RBACGrpcServer {
	return &RBACGrpcServer{
		roleCtrl:     r.NewRoleController(roleUC),
		permCtrl:     p.NewPermissionController(permUC),
		userRoleCtrl: ur.NewUserRoleController(userRoleUC),
		userPermCtrl: up.NewUserPermissionController(userPermUC),
		rolePermCtrl: rp.NewRolePermissionController(rolePermUC),
	}
}

func (s *RBACGrpcServer) CreateRole(ctx context.Context, req *rbacv1.CreateRoleRequest) (*rbacv1.CreateRoleResponse, error) {
	return s.roleCtrl.CreateRole(ctx, req)
}

func (s *RBACGrpcServer) GetRoleByID(ctx context.Context, req *rbacv1.GetRoleByNameRequest) (*rbacv1.GetRoleByNameResponse, error) {
	return s.roleCtrl.GetRoleByName(ctx, req)
}

func (s *RBACGrpcServer) ListRoles(ctx context.Context, req *rbacv1.Empty) (*rbacv1.ListRolesResponse, error) {
	return s.roleCtrl.ListRoles(ctx, req)
}

func (s *RBACGrpcServer) UpdateRole(ctx context.Context, req *rbacv1.UpdateRoleRequest) (*rbacv1.UpdateRoleResponse, error) {
	return s.roleCtrl.UpdateRole(ctx, req)
}

func (s *RBACGrpcServer) DeleteRole(ctx context.Context, req *rbacv1.DeleteRoleRequest) (*rbacv1.DeleteRoleResponse, error) {
	return s.roleCtrl.DeleteRole(ctx, req)
}

func (s *RBACGrpcServer) AssignRolesToUser(ctx context.Context, req *rbacv1.AssignRolesToUserRequest) (*rbacv1.AssignRolesToUserResponse, error) {
	return s.userRoleCtrl.AssignRolesToUser(ctx, req)
}

func (s *RBACGrpcServer) CreatePermission(ctx context.Context, req *rbacv1.CreatePermissionRequest) (*rbacv1.CreatePermissionResponse, error) {
	return s.permCtrl.CreatePermission(ctx, req)
}

func (s *RBACGrpcServer) DeletePermission(ctx context.Context, req *rbacv1.DeletePermissionRequest) (*rbacv1.DeletePermissionResponse, error) {
	return s.permCtrl.DeletePermission(ctx, req)
}

func (s *RBACGrpcServer) AssignPermissionsToRole(ctx context.Context, req *rbacv1.AssignPermissionsToRoleRequest) (*rbacv1.AssignPermissionsToRoleResponse, error) {
	return s.rolePermCtrl.AssignPermissionsToRole(ctx, req)
}

func (s *RBACGrpcServer) AssignPermissionsToUser(ctx context.Context, req *rbacv1.AssignPermissionsToUserRequest) (*rbacv1.AssignPermissionsToUserResponse, error) {
	return s.userPermCtrl.AssignPermissionsToUser(ctx, req)
}

func (s *RBACGrpcServer) ListPermissionsForRole(ctx context.Context, req *rbacv1.ListPermissionsForRoleRequest) (*rbacv1.ListPermissionsForRoleResponse, error) {
	return s.rolePermCtrl.ListPermissionsForRole(ctx, req)
}

func (s *RBACGrpcServer) ListPermissions(ctx context.Context, req *rbacv1.Empty) (*rbacv1.ListPermissionsResponse, error) {
	return s.permCtrl.ListPermissions(ctx, req)
}

func (s *RBACGrpcServer) ListPermissionsForUser(ctx context.Context, req *rbacv1.ListPermissionsForUserRequest) (*rbacv1.ListPermissionsForUserResponse, error) {
	return s.userPermCtrl.ListPermissionsForUser(ctx, req)
}

func (s *RBACGrpcServer) ListRolesForUser(ctx context.Context, req *rbacv1.ListRolesForUserRequest) (*rbacv1.ListRolesForUserResponse, error) {
	return s.userRoleCtrl.ListRolesForUser(ctx, req)
}

func (s *RBACGrpcServer) RemovePermissionFromRole(ctx context.Context, req *rbacv1.RemovePermissionFromRoleRequest) (*rbacv1.RemovePermissionFromRoleResponse, error) {
	return s.rolePermCtrl.RemovePermissionFromRole(ctx, req)
}

func (s *RBACGrpcServer) RemovePermissionFromUser(ctx context.Context, req *rbacv1.RemovePermissionFromUserRequest) (*rbacv1.RemovePermissionFromUserResponse, error) {
	return s.userPermCtrl.RemovePermissionFromUser(ctx, req)
}

func (s *RBACGrpcServer) RemoveRoleFromUser(ctx context.Context, req *rbacv1.RemoveRoleFromUserRequest) (*rbacv1.RemoveRoleFromUserResponse, error) {
	return s.userRoleCtrl.RemoveRoleFromUser(ctx, req)
}
