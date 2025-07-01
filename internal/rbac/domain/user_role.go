package domain

import "time"

type UserRole struct {
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
