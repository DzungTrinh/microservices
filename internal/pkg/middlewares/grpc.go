package middlewares

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/pkg/auth"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
)

func JWTVerifyInterceptor(ctx context.Context, req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	claims, err := auth.ExtractClaimsFromContext(ctx)
	if err != nil {
		logger.GetInstance().Errorf("Invalid token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	newCtx := context.WithValue(ctx, "claims", claims)
	newCtx = context.WithValue(newCtx, "user_id", claims.ID)

	return handler(newCtx, req)
}

func AdminOnlyInterceptor(ctx context.Context, req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	claims, err := auth.ExtractClaimsFromContext(ctx)
	if err != nil {
		logger.GetInstance().Errorf("Invalid token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	hasAdmin := false
	for _, role := range claims.Roles {
		if role == constants.RoleAdmin {
			hasAdmin = true
			break
		}
	}

	if !hasAdmin {
		logger.GetInstance().Warnf("Unauthorized: admin role required")
		return nil, status.Errorf(codes.PermissionDenied, "admin role required")
	}

	return handler(ctx, req)
}
