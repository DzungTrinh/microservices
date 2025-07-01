package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/internal/rbac/infras/mysql"
)

type userRoleRepository struct {
	db *sql.DB
	*mysql.Queries
}

func NewUserRoleRepository(db *sql.DB) repo.UserRoleRepository {
	return &userRoleRepository{
		db:      db,
		Queries: mysql.New(db),
	}
}

func (r *userRoleRepository) AssignRolesToUser(ctx context.Context, userRole domain.UserRole) error {
	return r.Queries.AssignRolesToUser(ctx, mysql.AssignRolesToUserParams{
		UserID: userRole.UserID,
		RoleID: userRole.RoleID,
		ID:     userRole.RoleID,
	})
}

func (r *userRoleRepository) ListRolesForUser(ctx context.Context, userID string) ([]domain.Role, error) {
	results, err := r.Queries.ListRolesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	roles := make([]domain.Role, len(results))
	for i, r := range results {
		roles[i] = domain.Role{
			ID:        r.ID,
			Name:      r.Name,
			BuiltIn:   r.BuiltIn,
			CreatedAt: r.CreatedAt,
			DeletedAt: r.DeletedAt,
		}
	}
	return roles, nil
}
