package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(50);not null" validate:"required,min=2,max=50"`
	Email    string `json:"email" gorm:"type:varchar(100);unique;not null" validate:"required,email"`
	Password string `json:"password" gorm:"type:varchar(255);not null" validate:"required,min=6"`
	Role     string `json:"role" gorm:"type:varchar(20);default:'user'"`
}

// BeforeCreate - Mã hóa password trước khi lưu vào database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
