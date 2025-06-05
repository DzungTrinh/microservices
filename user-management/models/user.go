package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role defines user roles
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email" gorm:"unique"`
	Password string `json:"password" validate:"required,min=6"`
	Role     Role   `json:"role" validate:"required"`
}

// BeforeCreate hook to hash password before saving to database
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
