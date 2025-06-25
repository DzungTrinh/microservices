package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/internal/rbac/infras/mysql"
)

type permissionRepository struct {
	db *sql.DB
	*mysql.Queries
}

func NewPermissionRepository(db *sql.DB) repo.PermissionRepository {
	return &permissionRepository{
		db:      db,
		Queries: mysql.New(db),
	}
}

func (r *permissionRepository) CreatePermission(ctx context.Context, perm domain.Permission) error {
	_, err := r.Queries.CreatePermission(ctx, mysql.CreatePermissionParams{
		ID:   perm.ID,
		Name: perm.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) DeletePermission(ctx context.Context, id string) error {
	return r.Queries.DeletePermission(ctx, id)
}

func (r *permissionRepository) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	results, err := r.Queries.ListPermissions(ctx)
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
