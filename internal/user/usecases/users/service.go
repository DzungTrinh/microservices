package users

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"microservices/user-management/internal/user/domain"
	"time"
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
	user, err := u.repo.CreateUser(ctx, req.Username, req.Email, string(hashedPassword))
	if err != nil {
		return domain.UserResp{}, err
	}
	return domain.UserResp{ID: user.ID, Username: user.Username, Email: user.Email}, nil
}

func (u *UserUsecase) Login(ctx context.Context, req domain.LoginReq) (domain.LoginResp, error) {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return domain.LoginResp{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return domain.LoginResp{}, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return domain.LoginResp{}, err
	}
	return domain.LoginResp{Token: tokenString}, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (domain.UserResp, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.UserResp{}, err
	}
	return domain.UserResp{ID: user.ID, Username: user.Username, Email: user.Email}, nil
}
