package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/user/usecases/auth"
	"microservices/user-management/pkg/logger"
	userv1 "microservices/user-management/proto/gen/user/v1"
)

type AuthController struct {
	uc     auth.AuthUseCase
	logger *logger.LoggerService
}

func NewAuthController(uc auth.AuthUseCase) *AuthController {
	return &AuthController{
		uc:     uc,
		logger: logger.GetInstance(),
	}
}

func (c *AuthController) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	userAgent := ctx.Value("user-agent").(string)
	ipAddress := ctx.Value("ip-address").(string)

	user, accessToken, refreshToken, err := c.uc.Register(ctx, req.Email, req.Username, req.Password, userAgent, ipAddress)
	if err != nil {
		c.logger.Errorf("Failed to register user: %v", err)
		return &userv1.RegisterResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &userv1.RegisterResponse{
		UserId:       user.ID,
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Success:      true,
	}, nil
}

func (c *AuthController) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	userAgent := ctx.Value("user-agent").(string)
	ipAddress := ctx.Value("ip-address").(string)

	user, accessToken, refreshToken, err := c.uc.Login(ctx, req.Email, req.Password, userAgent, ipAddress)
	if err != nil {
		c.logger.Errorf("Failed to login user: %v", err)
		code := codes.Unauthenticated
		if err.Error() != "invalid email or password" {
			code = codes.Internal
		}
		return &userv1.LoginResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(code, err.Error())
	}

	return &userv1.LoginResponse{
		UserId:       user.ID,
		Email:        user.Email,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Success:      true,
	}, nil
}
