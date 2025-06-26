package refresh_token

import (
	"context"
)

type RefreshTokenUseCase interface {
	RefreshToken(ctx context.Context, refreshToken, userAgent, ipAddress string) (string, string, error)
}
