package repo

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
}
