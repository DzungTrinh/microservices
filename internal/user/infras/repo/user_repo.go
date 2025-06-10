package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/infras/mysql"
)

type UserRepository struct {
	queries *mysql.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{queries: mysql.New(db)}
}

func (r *UserRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (domain.User, error) {
	arg := mysql.CreateUserParams{Username: username, Email: email, PasswordHash: passwordHash}
	result, err := r.queries.CreateUser(ctx, arg)
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
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}
