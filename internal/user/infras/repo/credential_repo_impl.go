package repo

import (
	"context"
	"database/sql"
	"microservices/user-management/internal/user/infras/mysql"

	"microservices/user-management/pkg/logger"

	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/domain/repo"
)

type credentialRepository struct {
	db      *sql.DB
	Queries *mysql.Queries
	logger  *logger.LoggerService
}

func NewCredentialRepository(db *sql.DB) repo.CredentialRepository {
	return &credentialRepository{
		db:      db,
		Queries: mysql.New(db),
		logger:  logger.GetInstance(),
	}
}

func (r *credentialRepository) CreateCredential(ctx context.Context, credential domain.Credential) error {
	err := r.Queries.CreateCredential(ctx, mysql.CreateCredentialParams{
		ID:          credential.ID,
		UserID:      credential.UserID,
		Provider:    credential.Provider,
		SecretHash:  credential.SecretHash,
		ProviderUid: credential.ProviderUID,
	})
	if err != nil {
		r.logger.Errorf("Failed to create credential for user %s: %v", credential.UserID, err)
		return err
	}
	return nil
}

func (r *credentialRepository) GetCredentialByEmailAndProvider(ctx context.Context, email, provider string) (domain.Credential, error) {
	credential, err := r.Queries.GetCredentialByEmailAndProvider(ctx, mysql.GetCredentialByEmailAndProviderParams{
		Email:    email,
		Provider: provider,
	})
	if err == sql.ErrNoRows {
		return domain.Credential{}, sql.ErrNoRows
	}
	if err != nil {
		r.logger.Errorf("Failed to get credential for email %s and provider %s: %v", email, provider, err)
		return domain.Credential{}, err
	}
	return domain.Credential{
		ID:          credential.ID,
		UserID:      credential.UserID,
		Provider:    credential.Provider,
		SecretHash:  credential.SecretHash,
		ProviderUID: credential.ProviderUid,
		CreatedAt:   credential.CreatedAt,
		DeletedAt:   credential.DeletedAt,
	}, nil
}
