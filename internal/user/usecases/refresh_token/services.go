package refresh_token

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"microservices/user-management/internal/pkg/auth"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/domain/repo"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
	"time"
)

type refreshTokenUseCase struct {
	rtRepo   repo.RefreshTokenRepository
	userRepo repo.UserRepository
}

func NewRefreshTokenUseCase(rtRepo repo.RefreshTokenRepository, userRepo repo.UserRepository) RefreshTokenUseCase {
	return &refreshTokenUseCase{rtRepo: rtRepo, userRepo: userRepo}
}

func (s *refreshTokenUseCase) RefreshToken(ctx context.Context, refreshToken, userAgent, ipAddress string) (string, string, error) {
	token, err := s.rtRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.GetInstance().Errorf("Failed to get refresh token: %v", err)
		return "", "", errors.New("invalid refresh token")
	}

	if token.Revoked || token.ExpiresAt.Before(time.Now()) {
		logger.GetInstance().Errorf("Refresh token %s is revoked or expired", token.ID)
		return "", "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetUserById(ctx, token.UserID)
	if err != nil {
		logger.GetInstance().Errorf("Failed to get user for refresh token %s: %v", token.ID, err)
		return "", "", errors.New("user not found")
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID, []string{constants.RoleUser}, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		logger.GetInstance().Errorf("Failed to generate token pair for user %s: %v", user.ID, err)
		return "", "", err
	}

	newRefreshTokenEntity := domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    token.UserID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	err = s.rtRepo.CreateRefreshToken(ctx, newRefreshTokenEntity)
	if err != nil {
		logger.GetInstance().Errorf("Failed to create new refresh token for user %s: %v", token.UserID, err)
		return "", "", err
	}

	err = s.rtRepo.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.GetInstance().Errorf("Failed to revoke old refresh token %s: %v", token.ID, err)
	}

	logger.GetInstance().Infof("Tokens refreshed for user %s", token.UserID)
	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}

func (u *refreshTokenUseCase) CleanExpiredTokens(ctx context.Context) error {
	return u.rtRepo.DeleteExpiredRefreshTokens(ctx)
}
