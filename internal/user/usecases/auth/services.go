package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"microservices/user-management/internal/pkg/auth"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/domain/repo"
	"microservices/user-management/internal/user/dto"
	"microservices/user-management/internal/user/infras/hash"
	"microservices/user-management/pkg/logger"
	"time"
)

type authUseCase struct {
	userRepo repo.UserRepository
	credRepo repo.CredentialRepository
	rtRepo   repo.RefreshTokenRepository
}

func NewAuthUseCase(userRepo repo.UserRepository, credRepo repo.CredentialRepository, rtRepo repo.RefreshTokenRepository) AuthUseCase {
	return &authUseCase{userRepo: userRepo, credRepo: credRepo, rtRepo: rtRepo}
}

func (s *authUseCase) Register(ctx context.Context, email, username, password, userAgent, ipAddress string) (domain.User, string, string, error) {
	if email == "" || username == "" || password == "" {
		return domain.User{}, "", "", errors.New("invalid email, username, or password")
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		logger.GetInstance().Errorf("Failed to hash password for user %s: %v", username, err)
		return domain.User{}, "", "", err
	}

	userID := uuid.New().String()
	user := domain.User{
		ID:            userID,
		Email:         email,
		Username:      username,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		logger.GetInstance().Errorf("Failed to create user %s: %v", userID, err)
		return domain.User{}, "", "", err
	}

	credential := domain.Credential{
		ID:         uuid.New().String(),
		UserID:     userID,
		Provider:   dto.ProviderLocal,
		SecretHash: hashedPassword,
		CreatedAt:  time.Now(),
	}
	err = s.credRepo.CreateCredential(ctx, credential)
	if err != nil {
		logger.GetInstance().Errorf("Failed to create credential for user %s: %v", userID, err)
		return domain.User{}, "", "", err
	}

	tokenPair, err := auth.GenerateTokenPair(userID, []string{dto.RoleUser}, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		logger.GetInstance().Errorf("Failed to generate token pair for user %s: %v", userID, err)
		return domain.User{}, "", "", err
	}

	refreshTokenEntity := domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	err = s.rtRepo.CreateRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		logger.GetInstance().Errorf("Failed to create refresh token for user %s: %v", userID, err)
		return domain.User{}, "", "", err
	}

	logger.GetInstance().Infof("User registered: %s (ID: %s)", username, userID)
	return user, tokenPair.AccessToken, tokenPair.RefreshToken, nil
}

func (s *authUseCase) Login(ctx context.Context, email, password, userAgent, ipAddress string) (domain.User, string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.GetInstance().Errorf("Failed to get user by email %s: %v", email, err)
		return domain.User{}, "", "", errors.New("invalid email or password")
	}

	credential, err := s.credRepo.GetCredentialByEmailAndProvider(ctx, email, dto.ProviderLocal)
	if err != nil {
		logger.GetInstance().Errorf("Failed to get credential for email %s: %v", email, err)
		return domain.User{}, "", "", errors.New("invalid email or password")
	}

	if err := hash.ComparePassword(credential.SecretHash, password); err != nil {
		logger.GetInstance().Errorf("Password mismatch for user %s", user.ID)
		return domain.User{}, "", "", errors.New("invalid email or password")
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID, []string{dto.RoleUser}, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		logger.GetInstance().Errorf("Failed to generate token pair for user %s: %v", user.ID, err)
		return domain.User{}, "", "", err
	}

	refreshTokenEntity := domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	err = s.rtRepo.CreateRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		logger.GetInstance().Errorf("Failed to create refresh token for user %s: %v", user.ID, err)
		return domain.User{}, "", "", err
	}

	logger.GetInstance().Infof("User logged in: %s (ID: %s)", user.Username, user.ID)
	return user, tokenPair.AccessToken, tokenPair.RefreshToken, nil
}
