package domain

import "time"

type Role struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BuiltIn   bool      `json:"built_in"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
