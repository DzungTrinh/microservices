package domain

import (
	"time"
)

type RegisterUserReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateUserModel struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Roles    []Role `json:"roles" binding:"required"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    []Role `json:"roles"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	MfaRequired  bool   `json:"mfa_required"`
}

type CreateRefreshTokenModel struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	UserAgent string    `json:"user_agent"`
	IpAddress string    `json:"ip_address"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RefreshTokenModel struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
