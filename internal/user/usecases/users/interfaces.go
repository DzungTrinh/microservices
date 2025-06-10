package users

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type (
	UseCase interface {
		Register(ctx context.Context, req domain.RegisterUserReq) (domain.UserResp, error)
		Login(ctx context.Context, req domain.LoginReq) (domain.LoginResp, error)
		GetUserByID(ctx context.Context, id int64) (domain.UserResp, error)
	}
)
