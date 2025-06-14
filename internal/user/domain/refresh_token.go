package domain

import (
	"time"
)

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	UserAgent string
	IpAddress string
	CreatedAt time.Time
	ExpiresAt time.Time
	Revoked   bool
}
