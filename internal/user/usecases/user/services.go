package user

import (
	"microservices/user-management/internal/user/domain/repo"
	"time"
)

type userUsecase struct {
	userRepo        repo.UserRepository
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}
