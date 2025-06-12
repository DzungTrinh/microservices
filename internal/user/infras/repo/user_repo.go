package repo

import (
	"context"
	"database/sql"
	"fmt"
	"microservices/user-management/internal/user/infras/mysql"

	"microservices/user-management/internal/user/domain"
)

type userRepository struct {
	queries *mysql.Queries
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{queries: mysql.New(db)}
}

func (r *userRepository) CreateUser(ctx context.Context, username, email, password string, role domain.Role) (domain.User, error) {
	if !domain.IsValidRole(role.String()) {
		return domain.User{}, fmt.Errorf("invalid role: %s", role)
	}

	result, err := r.queries.CreateUser(ctx, mysql.CreateUserParams{
		Username: username,
		Email:    email,
		Password: password,
		Role:     string(role),
	})
	if err != nil {
		return domain.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}

	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      domain.Role(user.Role),
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

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      domain.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		return domain.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      domain.Role(user.Role),
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
		result[i] = domain.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Password:  user.Password,
			Role:      domain.Role(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return result, nil
}
