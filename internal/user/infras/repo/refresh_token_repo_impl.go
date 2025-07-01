package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/domain/repo"
	"microservices/user-management/internal/user/infras/mysql"
	"microservices/user-management/pkg/logger"
)

type refreshTokenRepository struct {
	db      *sql.DB
	Queries *mysql.Queries
	logger  *logger.LoggerService
}

func NewRefreshTokenRepository(db *sql.DB) repo.RefreshTokenRepository {
	return &refreshTokenRepository{
		db:      db,
		Queries: mysql.New(db),
		logger:  logger.GetInstance(),
	}
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error {
	err := r.Queries.CreateRefreshToken(ctx, mysql.CreateRefreshTokenParams{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		Token:     refreshToken.Token,
		UserAgent: refreshToken.UserAgent,
		IpAddress: refreshToken.IPAddress,
		ExpiresAt: refreshToken.ExpiresAt,
		Revoked:   refreshToken.Revoked,
	})
	if err != nil {
		r.logger.Errorf("Failed to create refresh token for user %s: %v", refreshToken.UserID, err)
		return err
	}
	return nil
}

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, token string) (domain.RefreshToken, error) {
	refreshToken, err := r.Queries.GetRefreshToken(ctx, token)
	if err == sql.ErrNoRows {
		return domain.RefreshToken{}, sql.ErrNoRows
	}
	if err != nil {
		r.logger.Errorf("Failed to get refresh token %s: %v", token, err)
		return domain.RefreshToken{}, err
	}
	return domain.RefreshToken{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		Token:     refreshToken.Token,
		UserAgent: refreshToken.UserAgent,
		IPAddress: refreshToken.IpAddress,
		CreatedAt: refreshToken.CreatedAt,
		ExpiresAt: refreshToken.ExpiresAt,
		Revoked:   refreshToken.Revoked,
		DeletedAt: refreshToken.DeletedAt,
	}, nil
}

func (r *refreshTokenRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	err := r.Queries.RevokeRefreshToken(ctx, token)
	if err != nil {
		r.logger.Errorf("Failed to revoke refresh token %s: %v", token, err)
		return err
	}
	return nil
}

func (r *refreshTokenRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	err := r.Queries.DeleteExpiredRefreshTokens(ctx)
	if err != nil {
		r.logger.Errorf("Failed to delete expired tokens: %v", err)
		return err
	}
	return nil
}
