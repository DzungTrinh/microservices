package domain

import (
	"time"
)

// User represents the core entity.
type User struct {
	ID           int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string
	Email        string
	PasswordHash string
}
