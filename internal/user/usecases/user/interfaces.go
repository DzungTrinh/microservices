package user

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type UserUseCase interface {
	CreateAdmin(ctx context.Context, email, username, password string) (domain.User, error)
}
