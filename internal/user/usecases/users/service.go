package users

import (
	"context"
	"errors"
	"microservices/user-management/internal/pkg/auth"
	"time"

	"golang.org/x/crypto/bcrypt"
	"microservices/user-management/internal/user/domain"
)

type userUsecase struct {
	userRepo        domain.UserRepository
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserUsecase(userRepo domain.UserRepository) UserUseCase {
	return &userUsecase{
		userRepo:        userRepo,
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 7 * 24 * time.Hour,
	}
}

func (u *userUsecase) Register(ctx context.Context, req domain.RegisterUserReq) (domain.UserResp, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.UserResp{}, err
	}

	roles := make([]domain.Role, len(req.Roles))
	for i, role := range req.Roles {
		roles[i] = domain.Role(role)
	}
	if len(roles) == 0 {
		roles = []domain.Role{domain.RoleUser}
	}

	user, err := u.userRepo.CreateUser(ctx, req.Username, req.Email, string(hashedPassword), roles)
	if err != nil {
		return domain.UserResp{}, err
	}

	return domain.UserResp{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (map[string]string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	roles, err := u.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	role := "user"
	if len(roles) > 0 {
		role = roles[0].Name
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Email, role, "access", u.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateToken(user.ID, user.Email, "", "refresh", u.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (u *userUsecase) RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error) {
	claims, err := auth.VerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	roles, err := u.userRepo.GetUserRoles(ctx, claims.ID)
	if err != nil {
		return nil, err
	}

	role := "user"
	if len(roles) > 0 {
		role = roles[0].Name
	}

	accessToken, err := auth.GenerateToken(claims.ID, claims.Email, role, "access", u.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := auth.GenerateToken(claims.ID, claims.Email, "", "refresh", u.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (domain.UserResp, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return domain.UserResp{}, err
	}

	return domain.UserResp{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
	}, nil
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]domain.UserResp, error) {
	users, err := u.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.UserResp, len(users))
	for i, user := range users {
		result[i] = domain.UserResp{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.Roles,
		}
	}

	return result, nil
}

func (u *userUsecase) GetCurrentUser(ctx context.Context, userID string) (domain.UserResp, error) {
	return u.GetUserByID(ctx, userID)
}

func (u *userUsecase) UpdateUserRoles(ctx context.Context, userID string, roles []string) (domain.UserResp, error) {
	domainRoles := make([]domain.Role, len(roles))
	for i, role := range roles {
		domainRoles[i] = domain.Role(role)
	}

	if err := u.userRepo.UpdateUserRoles(ctx, userID, domainRoles); err != nil {
		return domain.UserResp{}, err
	}

	return u.GetUserByID(ctx, userID)
}
