package users

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, req domain.RegisterUserReq) (domain.AuthTokens, error)
	Login(ctx context.Context, req domain.LoginReq) (domain.AuthTokens, error)
	GetUserByID(ctx context.Context, id string) (domain.UserDTO, error)
	GetAllUsers(ctx context.Context) ([]domain.UserDTO, error)
	GetCurrentUser(ctx context.Context, userID string) (domain.UserDTO, error)
	UpdateUserRoles(ctx context.Context, userID string, roles []string) (domain.UserDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.AuthTokens, error)
	CleanExpiredTokens(ctx context.Context) error
}
