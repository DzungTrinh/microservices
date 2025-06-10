package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/pkg/errors"
)

func RequireRole(claims *Claims, role string) error {
	if claims.Role != role {
		return fmt.Errorf("requires %s role, got %s", role, claims.Role)
	}
	return nil
}

func RequireRoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			errors.HandleError(c, 401, errors.ErrUnauthorized, "Claims missing")
			c.Abort()
			return
		}

		if err := RequireRole(claims.(*Claims), role); err != nil {
			errors.HandleError(c, 403, errors.ErrForbidden, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRoleInterceptor(role string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims, ok := ctx.Value("claims").(*Claims)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "claims missing")
		}

		if err := RequireRole(claims, role); err != nil {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}

		return handler(ctx, req)
	}
}
