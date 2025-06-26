package domain

import "time"

type Credential struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Provider    string     `json:"provider"`
	SecretHash  string     `json:"secret_hash"`
	ProviderUID string     `json:"provider_uid"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
