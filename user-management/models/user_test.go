package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUser_BeforeCreate(t *testing.T) {
	db, err := gorm.Open(nil, &gorm.Config{}) // Mock GORM DB
	assert.NoError(t, err)

	user := User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     RoleUser,
	}

	err = user.BeforeCreate(db)
	assert.NoError(t, err)

	// Verify password is hashed
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(t, err)
}
