package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/internal/rbac/infras/mysql"
	"time"
)

type userPermissionRepository struct {
	db *sql.DB
	*mysql.Queries
}

func NewUserPermissionRepository(db *sql.DB) repo.UserPermissionRepository {
	return &userPermissionRepository{
		db:      db,
		Queries: mysql.New(db),
	}
}

func (r *userPermissionRepository) AssignPermissionsToUser(ctx context.Context, userPerm domain.UserPermission) error {
	expiresAt := userPerm.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = time.Time{}
	}
	granterID := userPerm.GranterID
	if granterID == "" {
		granterID = ""
	}
	return r.Queries.AssignPermissionsToUser(ctx, mysql.AssignPermissionsToUserParams{
		UserID:       userPerm.UserID,
		PermissionID: userPerm.PermissionID,
		GranterID:    granterID,
		ExpiresAt:    expiresAt,
		ID:           userPerm.PermissionID,
		GranterID_2:  granterID,
		ExpiresAt_2:  expiresAt,
	})
}

func (r *userPermissionRepository) ListPermissionsForUser(ctx context.Context, userID string) ([]domain.Permission, error) {
	results, err := r.Queries.ListPermissionsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	permissions := make([]domain.Permission, len(results))
	for i, p := range results {
		permissions[i] = domain.Permission{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			DeletedAt: p.DeletedAt,
		}
	}
	return permissions, nil
}
