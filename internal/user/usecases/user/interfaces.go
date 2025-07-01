package user

import (
	"context"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/dto"
)

type UserUseCase interface {
	CreateAdmin(ctx context.Context, email, username, password string) (domain.User, error)
	GetAllUsers(ctx context.Context) ([]dto.UserDTO, error)
}
