package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"microservices/user-management/internal/pkg/auth"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/domain/repo"
	"microservices/user-management/internal/user/dto"
	"microservices/user-management/internal/user/infras/hash"
	"microservices/user-management/internal/user/usecases/grpc"
	"microservices/user-management/pkg/logger"
	"time"
)

type authUseCase struct {
	userRepo   repo.UserRepository
	credRepo   repo.CredentialRepository
	rtRepo     repo.RefreshTokenRepository
	outboxRepo repo.OutboxRepository
	txManager  repo.TxManager
	rbacClient grpc.RBACService
	accessTtl  time.Duration
	refreshTtl time.Duration
}

func NewAuthUseCase(userRepo repo.UserRepository, credRepo repo.CredentialRepository, rtRepo repo.RefreshTokenRepository,
	outboxRepo repo.OutboxRepository, txManager repo.TxManager, rbacClient grpc.RBACService) AuthUseCase {
	return &authUseCase{userRepo: userRepo, credRepo: credRepo, rtRepo: rtRepo,
		outboxRepo: outboxRepo, txManager: txManager,
		accessTtl: time.Minute * 15, refreshTtl: time.Hour * 24 * 7,
		rbacClient: rbacClient}
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

	credential := domain.Credential{
		ID:          uuid.New().String(),
		UserID:      userID,
		Provider:    dto.ProviderLocal,
		ProviderUID: userID,
		SecretHash:  hashedPassword,
		CreatedAt:   time.Now(),
	}

	var accessToken, refreshToken string

	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.CreateUser(txCtx, user); err != nil {
			return err
		}
		if err := s.credRepo.CreateCredential(txCtx, credential); err != nil {
			return err
		}

		// build outbox event
		rolePayload := struct {
			UserID string `json:"user_id"`
		}{
			UserID: userID,
		}

		payloadBytes, err := json.Marshal(rolePayload)
		if err != nil {
			return err
		}

		outboxEvent := domain.OutboxEvent{
			AggregateType: "User",
			AggregateID:   userID,
			Type:          "UserRegistered",
			Payload:       string(payloadBytes),
			Status:        dto.OutboxPending,
			CreatedAt:     time.Now(),
		}

		if err := s.outboxRepo.InsertEvent(txCtx, &outboxEvent); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.GetInstance().Errorf("Failed to register user %s: %v", userID, err)
		return domain.User{}, "", "", err
	}

	// Call Login to generate tokens
	createdUser, accessToken, refreshToken, err := s.Login(ctx, email, password, userAgent, ipAddress)
	if err != nil {
		logger.GetInstance().Errorf("Failed to login after registering user %s: %v", userID, err)
		return createdUser, "", "", err
	}

	logger.GetInstance().Infof("User registered: %s (ID: %s)", username, userID)
	return user, accessToken, refreshToken, nil
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

	// Fetch user roles from RBAC service
	roles, err := s.rbacClient.ListRolesForUser(ctx, user.ID)
	if err != nil {
		logger.GetInstance().Errorf("Failed to fetch roles for user %s: %v", user.ID, err)
		return domain.User{}, "", "", err
	}

	// Fetch user permissions from RBAC service
	perms, err := s.rbacClient.ListPermissionsForUser(ctx, user.ID)
	if err != nil {
		logger.GetInstance().Errorf("Failed to fetch permissions for user %s: %v", user.ID, err)
		return domain.User{}, "", "", err
	}

	// Generate tokens with actual roles
	tokenPair, err := auth.GenerateTokenPair(user.ID, roles, perms, s.accessTtl, s.refreshTtl)
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
		ExpiresAt: time.Now().Add(s.refreshTtl),
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
