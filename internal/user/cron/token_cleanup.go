package cron

import (
	"context"
	"microservices/user-management/internal/user/usecases/refresh_token"
	"microservices/user-management/pkg/logger"
	"time"
)

// StartTokenCleanup runs a periodic task to clean expired tokens
func StartTokenCleanup(ctx context.Context, usecase refresh_token.RefreshTokenUseCase) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := usecase.CleanExpiredTokens(ctx); err != nil {
					logger.GetInstance().Printf("Failed to clean expired tokens: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
