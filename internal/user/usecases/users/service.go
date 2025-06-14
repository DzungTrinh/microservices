package users

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"microservices/user-management/internal/pkg/auth"
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
	req.Password = string(hashedPassword)

	user, err := u.userRepo.CreateUser(ctx, req)
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

func (u *userUsecase) Login(ctx context.Context, req domain.LoginReq) (domain.LoginResp, error) {
	log.Printf("Login attempt for email: %s", req.Email)
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("GetUserByEmail failed: %v", err)
		return domain.LoginResp{}, errors.New("invalid credentials")
	}
	log.Printf("Found user: %+v", user)
	log.Printf("Stored hash: %s, Input password: %s", user.Password, req.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return domain.LoginResp{}, errors.New("invalid credentials")
	}

	roles := user.Roles
	role := string(domain.RoleUser)
	if len(roles) > 0 {
		role = string(roles[0])
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID, role, u.accessTokenTTL, u.refreshTokenTTL)
	if err != nil {
		log.Printf("Generate token pair failed: %v", err)
		return domain.LoginResp{}, err
	}

	refreshTokenID := uuid.New().String()
	userAgent, _ := ctx.Value("user-agent").(string)
	ipAddress, _ := ctx.Value("ip-address").(string)
	log.Printf("Context: User-Agent=%s, IP=%s", userAgent, ipAddress)

	if err := u.userRepo.CreateRefreshToken(ctx, domain.CreateRefreshTokenModel{
		ID:        refreshTokenID,
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
	}); err != nil {
		log.Printf("Error storing refresh token: %v", err)
		return domain.LoginResp{}, err
	}

	return domain.LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (u *userUsecase) RefreshToken(ctx context.Context, refreshToken string) (domain.LoginResp, error) {
	token, err := u.userRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil || token.Revoked {
		log.Printf("Invalid or revoked refresh token: %v", err)
		return domain.LoginResp{}, errors.New("invalid or revoked refresh token")
	}

	if time.Now().After(token.ExpiresAt) {
		log.Printf("Expired refresh token")
		return domain.LoginResp{}, errors.New("expired refresh token")
	}

	claims, err := auth.VerifyToken(refreshToken, "refresh")
	if err != nil {
		log.Printf("Invalid refresh token claims: %v", err)
		return domain.LoginResp{}, errors.New("invalid refresh token")
	}

	refreshClaims, ok := claims.(*auth.RefreshClaims)
	if !ok || refreshClaims.ID != token.UserID {
		log.Printf("Invalid refresh token claims type or ID mismatch")
		return domain.LoginResp{}, errors.New("invalid refresh token")
	}

	user, err := u.userRepo.GetUserByID(ctx, token.UserID)
	if err != nil {
		log.Printf("GetUserByID failed: %v", err)
		return domain.LoginResp{}, err
	}

	roles := user.Roles
	role := string(domain.RoleUser)
	if len(roles) > 0 {
		role = string(roles[0])
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID, role, u.accessTokenTTL, u.refreshTokenTTL)
	if err != nil {
		log.Printf("Generate token pair failed: %v", err)
		return domain.LoginResp{}, err
	}

	newRefreshTokenID := uuid.New().String()
	userAgent, _ := ctx.Value("user-agent").(string)
	ipAddress, _ := ctx.Value("ip-address").(string)
	log.Printf("Context: User-Agent=%s, IP=%s", userAgent, ipAddress)

	if err := u.userRepo.CreateRefreshToken(ctx, domain.CreateRefreshTokenModel{
		ID:        newRefreshTokenID,
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
	}); err != nil {
		log.Printf("Error storing refresh token: %v", err)
		return domain.LoginResp{}, err
	}

	if err := u.userRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		log.Printf("Revoke refresh token failed: %v", err)
		return domain.LoginResp{}, err
	}

	return domain.LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
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

func (u *userUsecase) CleanExpiredTokens(ctx context.Context) error {
	return u.userRepo.DeleteExpiredRefreshTokens(ctx)
}
