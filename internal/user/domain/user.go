package domain

import (
	"time"
)

type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Email     string
	Password  string
	Role      Role
}
