package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/domain/repo"
	"microservices/user-management/internal/rbac/infras/mysql"
)

type roleRepository struct {
	db *sql.DB
	*mysql.Queries
}

func NewRoleRepository(db *sql.DB) repo.RoleRepository {
	return &roleRepository{
		db:      db,
		Queries: mysql.New(db),
	}
}

func (r *roleRepository) CreateRole(ctx context.Context, role domain.Role) (domain.Role, error) {
	err := r.Queries.CreateRole(ctx, mysql.CreateRoleParams{
		ID:      role.ID,
		Name:    role.Name,
		BuiltIn: role.BuiltIn,
	})
	if err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id string) (domain.Role, error) {
	result, err := r.Queries.GetRoleByID(ctx, id)
	if err != nil {
		return domain.Role{}, err
	}
	return domain.Role{
		ID:        result.ID,
		Name:      result.Name,
		BuiltIn:   result.BuiltIn,
		CreatedAt: result.CreatedAt,
		DeletedAt: result.DeletedAt,
	}, nil
}

func (r *roleRepository) ListRoles(ctx context.Context) ([]domain.Role, error) {
	results, err := r.Queries.ListRoles(ctx)
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

func (r *roleRepository) UpdateRole(ctx context.Context, role domain.Role) (domain.Role, error) {
	err := r.Queries.UpdateRole(ctx, mysql.UpdateRoleParams{
		ID:      role.ID,
		Name:    role.Name,
		BuiltIn: role.BuiltIn,
	})
	if err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) DeleteRole(ctx context.Context, id string) error {
	return r.Queries.DeleteRole(ctx, id)
}
