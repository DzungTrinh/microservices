package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/domain/repo"
	"microservices/user-management/internal/user/dto"
	"microservices/user-management/internal/user/infras/hash"
	"microservices/user-management/pkg/logger"
	"time"
)

type userUseCase struct {
	userRepo   repo.UserRepository
	credRepo   repo.CredentialRepository
	outboxRepo repo.OutboxRepository
	txManager  repo.TxManager
}

func NewUserUseCase(userRepo repo.UserRepository, credRepo repo.CredentialRepository,
	outboxRepo repo.OutboxRepository, txManager repo.TxManager) UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		credRepo:   credRepo,
		outboxRepo: outboxRepo,
		txManager:  txManager,
	}
}

func (s *userUseCase) CreateAdmin(ctx context.Context, email, username, password string) (domain.User, error) {
	if email == "" || username == "" || password == "" {
		return domain.User{}, errors.New("invalid email, username, or password")
	}

	// Check if user exists
	if _, err := s.userRepo.GetUserByEmail(ctx, email); err == nil {
		return domain.User{}, errors.New("admin user already exists")
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		logger.GetInstance().Errorf("Failed to hash password for admin %s: %v", username, err)
		return domain.User{}, err
	}

	userID := uuid.New().String()
	user := domain.User{
		ID:            userID,
		Email:         email,
		Username:      username,
		EmailVerified: true, // Admin accounts are pre-verified
	}

	credential := domain.Credential{
		ID:          uuid.New().String(),
		UserID:      userID,
		Provider:    dto.ProviderLocal,
		ProviderUID: userID,
		SecretHash:  hashedPassword,
	}

	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.CreateUser(txCtx, user); err != nil {
			return err
		}
		if err := s.credRepo.CreateCredential(txCtx, credential); err != nil {
			return err
		}

		// Build outbox event
		payload := struct {
			UserID string `json:"user_id"`
		}{
			UserID: userID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		outboxEvent := domain.OutboxEvent{
			AggregateType: "User",
			AggregateID:   userID,
			Type:          "AdminUserCreated",
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
		logger.GetInstance().Errorf("Failed to create admin %s: %v", userID, err)
		return domain.User{}, err
	}

	logger.GetInstance().Infof("Admin user %s created: ID %s", username, userID)
	return user, nil
}
