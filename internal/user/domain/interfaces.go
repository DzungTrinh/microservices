package domain

import (
	"context"
)

// UserRepository defines the data access interface.
type UserRepository interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
}
