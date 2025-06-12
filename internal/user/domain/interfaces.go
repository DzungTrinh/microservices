package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, username, email, password string, role Role) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
}
