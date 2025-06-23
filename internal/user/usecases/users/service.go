package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
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

func (u *userUsecase) Register(ctx context.Context, req domain.RegisterUserReq) (domain.AuthTokens, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing failed: %v", err)
		return domain.AuthTokens{}, err
	}

	user, err := u.userRepo.CreateUser(ctx, req.Username, req.Email, string(hashedPassword), []domain.Role{domain.RoleUser})
	if err != nil {
		log.Printf("Create user failed: %v", err)
		if strings.Contains(err.Error(), "Duplicate entry") {
			return domain.AuthTokens{}, fmt.Errorf("email or username already exists")
		}
		return domain.AuthTokens{}, err
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID, string(domain.RoleUser), 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		log.Printf("Generate token pair failed: %v", err)
		return domain.AuthTokens{}, err
	}

	refreshTokenID := uuid.New().String()
	userAgent, _ := ctx.Value("user-agent").(string)
	ipAddress, _ := ctx.Value("ip-address").(string)

	if err := u.userRepo.CreateRefreshToken(ctx, domain.CreateRefreshTokenModel{
		ID:        refreshTokenID,
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
	}); err != nil {
		log.Printf("Error storing refresh token: %v", err)
		return domain.AuthTokens{}, err
	}

	return domain.AuthTokens{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		MfaRequired:  false,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, req domain.LoginReq) (domain.AuthTokens, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("GetUserByEmail failed: %v", err)
		return domain.AuthTokens{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return domain.AuthTokens{}, errors.New("invalid credentials")
	}

	roles := user.Roles
	role := string(domain.RoleUser)
	if len(roles) > 0 {
		role = string(roles[0])
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID, role, u.accessTokenTTL, u.refreshTokenTTL)
	if err != nil {
		log.Printf("Generate token pair failed: %v", err)
		return domain.AuthTokens{}, err
	}

	refreshTokenID := uuid.New().String()
	userAgent, _ := ctx.Value("user-agent").(string)
	ipAddress, _ := ctx.Value("ip-address").(string)

	if err := u.userRepo.CreateRefreshToken(ctx, domain.CreateRefreshTokenModel{
		ID:        refreshTokenID,
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
	}); err != nil {
		log.Printf("Error storing refresh token: %v", err)
		return domain.AuthTokens{}, err
	}

	return domain.AuthTokens{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		MfaRequired:  false,
	}, nil
}

func (u *userUsecase) RefreshToken(ctx context.Context, refreshToken string) (domain.AuthTokens, error) {
	token, err := u.userRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil || token.Revoked {
		log.Printf("Invalid or revoked refresh token: %v", err)
		return domain.AuthTokens{}, errors.New("invalid or revoked refresh token")
	}

	if time.Now().After(token.ExpiresAt) {
		log.Printf("Expired refresh token")
		return domain.AuthTokens{}, errors.New("expired refresh token")
	}

	claims, err := auth.VerifyToken(refreshToken, "refresh")
	if err != nil {
		log.Printf("Invalid refresh token claims: %v", err)
		return domain.AuthTokens{}, errors.New("invalid refresh token")
	}

	refreshClaims, ok := claims.(*auth.RefreshClaims)
	if !ok || refreshClaims.ID != token.UserID {
		log.Printf("Invalid refresh token claims type or ID mismatch")
		return domain.AuthTokens{}, errors.New("invalid refresh token")
	}

	user, err := u.userRepo.GetUserByID(ctx, token.UserID)
	if err != nil {
		log.Printf("GetUserByID failed: %v", err)
		return domain.AuthTokens{}, err
	}

	roles := user.Roles
	role := string(domain.RoleUser)
	if len(roles) > 0 {
		role = string(roles[0])
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID, role, u.accessTokenTTL, u.refreshTokenTTL)
	if err != nil {
		log.Printf("Generate token pair failed: %v", err)
		return domain.AuthTokens{}, err
	}

	newRefreshTokenID := uuid.New().String()
	userAgent, _ := ctx.Value("user-agent").(string)
	ipAddress, _ := ctx.Value("ip-address").(string)

	if err := u.userRepo.CreateRefreshToken(ctx, domain.CreateRefreshTokenModel{
		ID:        newRefreshTokenID,
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
	}); err != nil {
		log.Printf("Error storing refresh token: %v", err)
		return domain.AuthTokens{}, err
	}

	if err := u.userRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		log.Printf("Revoke refresh token failed: %v", err)
		return domain.AuthTokens{}, err
	}

	return domain.AuthTokens{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (domain.UserDTO, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return domain.UserDTO{}, err
	}

	return domain.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
	}, nil
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]domain.UserDTO, error) {
	users, err := u.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.UserDTO, len(users))
	for i, user := range users {
		result[i] = domain.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.Roles,
		}
	}

	return result, nil
}

func (u *userUsecase) GetCurrentUser(ctx context.Context, userID string) (domain.UserDTO, error) {
	return u.GetUserByID(ctx, userID)
}

func (u *userUsecase) UpdateUserRoles(ctx context.Context, userID string, roles []string) (domain.UserDTO, error) {
	domainRoles := make([]domain.Role, len(roles))
	for i, role := range roles {
		domainRoles[i] = domain.Role(role)
	}

	if err := u.userRepo.UpdateUserRoles(ctx, userID, domainRoles); err != nil {
		return domain.UserDTO{}, err
	}

	return u.GetUserByID(ctx, userID)
}

func (u *userUsecase) CreateUserAdmin(ctx context.Context, req domain.CreateUserModel) (domain.User, error) {
	if len(req.Roles) == 0 {
		log.Printf("No roles provided for admin user creation")
		return domain.User{}, fmt.Errorf("at least one role required")
	}
	for _, role := range req.Roles {
		if !domain.IsValidRole(string(role)) {
			log.Printf("Invalid role: %s", role)
			return domain.User{}, fmt.Errorf("invalid role: %s", role)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing failed: %v", err)
		return domain.User{}, err
	}

	user, err := u.userRepo.CreateUser(ctx, req.Username, req.Email, string(hashedPassword), req.Roles)
	if err != nil {
		log.Printf("Create admin user failed: %v", err)
		if strings.Contains(err.Error(), "Duplicate entry") {
			return domain.User{}, fmt.Errorf("email or username already exists")
		}
		return domain.User{}, err
	}

	return user, nil
}

func (u *userUsecase) CleanExpiredTokens(ctx context.Context) error {
	return u.userRepo.DeleteExpiredRefreshTokens(ctx)
}
