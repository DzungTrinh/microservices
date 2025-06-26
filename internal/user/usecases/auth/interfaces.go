package auth

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, username, password, userAgent, ipAddress string) (domain.User, string, string, error)
	Login(ctx context.Context, email, password, userAgent, ipAddress string) (domain.User, string, string, error)
}
