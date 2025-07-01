package grpc

import "context"

type RBACService interface {
	ListRolesForUser(ctx context.Context, userID string) ([]string, error)
	ListPermissionsForUser(ctx context.Context, userID string) ([]string, error)
}
