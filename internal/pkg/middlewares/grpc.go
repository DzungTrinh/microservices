package middlewares

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/pkg/auth"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
	"strings"
)

func ExtractClaimsFromContext(ctx context.Context) (*auth.AccessClaims, error) {
	// First, try getting from context
	if claims, ok := ctx.Value("claims").(*auth.AccessClaims); ok {
		return claims, nil
	}

	// Fallback: extract from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := auth.VerifyToken(token, "access")
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	accessClaims, ok := claims.(*auth.AccessClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return accessClaims, nil
}

func JWTVerifyInterceptor(ctx context.Context, req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	claims, err := ExtractClaimsFromContext(ctx)
	if err != nil {
		logger.GetInstance().Errorf("Invalid token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	newCtx := context.WithValue(ctx, "claims", claims)
	newCtx = context.WithValue(newCtx, "user_id", claims.ID)

	return handler(newCtx, req)
}

func AdminOnlyInterceptor(ctx context.Context, req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	claims, err := ExtractClaimsFromContext(ctx)
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
