package users

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"microservices/user-management/internal/pkg/auth"
	"microservices/user-management/internal/user/domain"
)

type UserUsecase struct {
	repo      domain.UserRepository
	jwtSecret []byte
}

func NewUserUsecase(repo domain.UserRepository, jwtSecret string) *UserUsecase {
	return &UserUsecase{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (u *UserUsecase) Register(ctx context.Context, req domain.RegisterUserReq) (domain.UserResp, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.UserResp{}, err
	}

	user, err := u.repo.CreateUser(ctx, req.Username, req.Email, string(hashedPassword), domain.RoleUser)
	if err != nil {
		return domain.UserResp{}, err
	}

	return domain.UserResp{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (u *UserUsecase) Login(ctx context.Context, req domain.LoginReq) (domain.LoginResp, error) {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return domain.LoginResp{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return domain.LoginResp{}, err
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, user.Role.String())
	if err != nil {
		return domain.LoginResp{}, err
	}

	return domain.LoginResp{Token: token}, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (domain.UserResp, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.UserResp{}, err
	}

	return domain.UserResp{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (u *UserUsecase) GetAllUsers(ctx context.Context) ([]domain.UserResp, error) {
	users, err := u.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.UserResp, len(users))
	for i, user := range users {
		result[i] = domain.UserResp{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}
	}

	return result, nil
}
