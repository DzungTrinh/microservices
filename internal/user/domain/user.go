package domain

import (
	"time"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	Roles     []Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
