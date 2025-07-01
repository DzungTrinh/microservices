package refresh_token

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/user/usecases/refresh_token"
	"microservices/user-management/pkg/logger"
	userv1 "microservices/user-management/proto/gen/user/v1"
)

type RefreshTokenController struct {
	uc     refresh_token.RefreshTokenUseCase
	logger *logger.LoggerService
}

func NewRefreshTokenController(uc refresh_token.RefreshTokenUseCase) *RefreshTokenController {
	return &RefreshTokenController{
		uc:     uc,
		logger: logger.GetInstance(),
	}
}

func (c *RefreshTokenController) RefreshToken(ctx context.Context, req *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	userAgent := ctx.Value("user-agent").(string)
	ipAddress := ctx.Value("ip-address").(string)

	accessToken, refreshToken, err := c.uc.RefreshToken(ctx, req.RefreshToken, userAgent, ipAddress)
	if err != nil {
		c.logger.Errorf("Failed to refresh token: %v", err)
		code := codes.Unauthenticated
		if err.Error() != "invalid refresh token" && err.Error() != "user not found" {
			code = codes.Internal
		}
		return &userv1.RefreshTokenResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(code, err.Error())
	}

	return &userv1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Success:      true,
	}, nil
}
