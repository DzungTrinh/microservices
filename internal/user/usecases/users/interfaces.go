package users

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type (
	UserUseCase interface {
		Register(ctx context.Context, req domain.RegisterUserReq) (domain.UserResp, error)
		Login(ctx context.Context, req domain.LoginReq) (domain.LoginResp, error)
		GetUserByID(ctx context.Context, id string) (domain.UserResp, error)
		GetAllUsers(ctx context.Context) ([]domain.UserResp, error)
		GetCurrentUser(ctx context.Context, userID string) (domain.UserResp, error)
		UpdateUserRoles(ctx context.Context, userID string, roles []string) (domain.UserResp, error)
		RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error)
	}
)
