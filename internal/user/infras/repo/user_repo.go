package repo

import (
	"context"
	"database/sql"
	"fmt"
	"microservices/user-management/internal/user/infras/mysql"
	"strings"

	"github.com/google/uuid"
	"microservices/user-management/internal/user/domain"
)

type userRepository struct {
	queries *mysql.Queries
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{queries: mysql.New(db)}
}

func (r *userRepository) CreateUser(ctx context.Context, username, email, password string, roles []domain.Role) (domain.User, error) {
	if len(roles) == 0 {
		roles = []domain.Role{domain.RoleUser}
	}
	for _, role := range roles {
		if !domain.IsValidRole(string(role)) {
			return domain.User{}, fmt.Errorf("invalid role: %s", role)
		}
	}

	id := uuid.New().String()
	_, err := r.queries.CreateUser(ctx, mysql.CreateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return domain.User{}, err
	}

	for _, role := range roles {
		roleID, err := r.queries.GetRoleIDByName(ctx, string(role))
		if err != nil {
			return domain.User{}, fmt.Errorf("role %s not found: %w", role, err)
		}
		if err := r.queries.CreateUserRole(ctx, mysql.CreateUserRoleParams{
			UserID: id,
			RoleID: roleID,
		}); err != nil {
			return domain.User{}, err
		}
	}

	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	var userRoles []domain.Role
	if user.Roles.Valid {
		for _, role := range strings.Split(user.Roles.String, ",") {
			userRoles = append(userRoles, domain.Role(role))
		}
	}

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     userRoles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return domain.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return domain.User{}, err
	}

	var userRoles []domain.Role
	if user.Roles.Valid {
		for _, role := range strings.Split(user.Roles.String, ",") {
			userRoles = append(userRoles, domain.Role(role))
		}
	}

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     userRoles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		return domain.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return domain.User{}, err
	}

	var userRoles []domain.Role
	if user.Roles.Valid {
		for _, role := range strings.Split(user.Roles.String, ",") {
			userRoles = append(userRoles, domain.Role(role))
		}
	}

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     userRoles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.User, len(users))
	for i, user := range users {
		var userRoles []domain.Role
		if user.Roles.Valid {
			for _, role := range strings.Split(user.Roles.String, ",") {
				userRoles = append(userRoles, domain.Role(role))
			}
		}
		result[i] = domain.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Password:  user.Password,
			Roles:     userRoles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return result, nil
}

func (r *userRepository) UpdateUserRoles(ctx context.Context, userID string, roles []domain.Role) error {
	if len(roles) == 0 {
		return fmt.Errorf("at least one role required")
	}
	for _, role := range roles {
		if !domain.IsValidRole(string(role)) {
			return fmt.Errorf("invalid role: %s", role)
		}
	}

	// Delete existing roles
	if err := r.queries.DeleteUserRoles(ctx, userID); err != nil {
		return err
	}

	// Insert new roles
	for _, role := range roles {
		roleID, err := r.queries.GetRoleIDByName(ctx, string(role))
		if err != nil {
			return fmt.Errorf("role %s not found: %w", role, err)
		}
		if err := r.queries.CreateUserRole(ctx, mysql.CreateUserRoleParams{
			UserID: userID,
			RoleID: roleID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	return r.queries.GetUserRoles(ctx, userID)
}
