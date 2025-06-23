package domain

import (
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, email, password string, roles []Role) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	UpdateUserRoles(ctx context.Context, userID string, roles []Role) error
	CreateRefreshToken(ctx context.Context, model CreateRefreshTokenModel) error
	GetRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
}
