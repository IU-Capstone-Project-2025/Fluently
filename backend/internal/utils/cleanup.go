package utils

import (
	"context"
	"time"

	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/pkg/logger"

	"go.uber.org/zap"
)

// CleanupExpiredTokens deletes expired link tokens
func CleanupExpiredTokens(repo *postgres.LinkTokenRepository) {
	ctx := context.Background()

	if err := repo.DeleteExpired(ctx); err != nil {
		logger.Log.Error("Failed to cleanup expired tokens", zap.Error(err))
	} else {
		logger.Log.Info("Successfully cleaned up expired link tokens")
	}
}

// StartTokenCleanupTask starts a periodic token cleanup task
func StartTokenCleanupTask(repo *postgres.LinkTokenRepository, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			CleanupExpiredTokens(repo)
		}
	}()

	logger.Log.Info("Started token cleanup task",
		zap.Duration("interval", interval))
}
