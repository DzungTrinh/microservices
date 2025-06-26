package repo

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserById(ctx context.Context, id string) (domain.User, error)
}
