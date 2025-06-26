package domain

import "time"

type RolePermission struct {
	RoleID       string    `json:"role_id"`
	PermissionID string    `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
