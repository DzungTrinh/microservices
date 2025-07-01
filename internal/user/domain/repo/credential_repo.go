package repo

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type CredentialRepository interface {
	CreateCredential(ctx context.Context, credential domain.Credential) error
	GetCredentialByEmailAndProvider(ctx context.Context, email, provider string) (domain.Credential, error)
}
