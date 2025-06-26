package router

import (
	"context"
	"google.golang.org/grpc"
	"microservices/user-management/internal/pkg/middlewares"
)

func InterceptorChain() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		switch info.FullMethod {
		case "/rbac.v1.RBACService/CreateRole",
			"/rbac.v1.RBACService/UpdateRole",
			"/rbac.v1.RBACService/DeleteRole",
			"/rbac.v1.RBACService/CreatePermission",
			"/rbac.v1.RBACService/DeletePermission",
			"/rbac.v1.RBACService/AssignRolesToUser",
			"/rbac.v1.RBACService/AssignPermissionsToRole",
			"/rbac.v1.RBACService/AssignPermissionsToUser":
			authCtx, err := middlewares.JWTVerifyInterceptor(ctx, req, func(c context.Context, _ interface{}) (interface{}, error) {
				return c, nil
			})
			if err != nil {
				return nil, err
			}
			return middlewares.AdminOnlyInterceptor(authCtx.(context.Context), req, handler)

		case "/rbac.v1.RBACService/GetRoleByID",
			"/rbac.v1.RBACService/ListRoles",
			"/rbac.v1.RBACService/ListPermissions",
			"/rbac.v1.RBACService/ListPermissionsForRole",
			"/rbac.v1.RBACService/ListPermissionsForUser",
			"/rbac.v1.RBACService/ListRolesForUser":
			return middlewares.JWTVerifyInterceptor(ctx, req, handler)

		default:
			return handler(ctx, req)
		}
	}
}
