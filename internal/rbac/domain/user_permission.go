package domain

import "time"

type UserPermission struct {
	UserID    string    `json:"user_id"`
	PermID    string    `json:"perm_id"`
	GranterID string    `json:"granter_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
