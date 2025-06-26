package domain

import (
	"time"
)

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}
