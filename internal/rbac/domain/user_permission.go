package domain

import "time"

type UserPermission struct {
	UserID       string    `json:"user_id"`
	PermissionID string    `json:"permission_id"`
	GranterID    string    `json:"granter_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
