package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/internal/rbac/infras/mysql"
)

type rolePermissionRepository struct {
	db *sql.DB
	*mysql.Queries
}

func NewRolePermissionRepository(db *sql.DB) repo.RolePermissionRepository {
	return &rolePermissionRepository{
		db:      db,
		Queries: mysql.New(db),
	}
}

func (r *rolePermissionRepository) AssignPermissionsToRole(ctx context.Context, rolePerm domain.RolePermission) error {
	return r.Queries.AssignPermissionsToRole(ctx, mysql.AssignPermissionsToRoleParams{
		RoleID:       rolePerm.RoleID,
		PermissionID: rolePerm.PermissionID,
	})
}

func (r *rolePermissionRepository) ListPermissionsForRole(ctx context.Context, roleID string) ([]domain.Permission, error) {
	results, err := r.Queries.ListPermissionsForRole(ctx, roleID)
	if err != nil {
		return nil, err
	}
	perms := make([]domain.Permission, len(results))
	for i, p := range results {
		perms[i] = domain.Permission{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			DeletedAt: p.DeletedAt,
		}
	}
	return perms, nil
}

func (r *rolePermissionRepository) SoftDeleteRolePermission(ctx context.Context, rolePerm domain.RolePermission) error {
	return r.Queries.SoftDeleteRolePermission(ctx, mysql.SoftDeleteRolePermissionParams{
		RoleID:       rolePerm.RoleID,
		PermissionID: rolePerm.PermissionID,
	})
}
