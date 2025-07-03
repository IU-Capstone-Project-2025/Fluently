// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fluently/telegram-bot/config"
	"fluently/telegram-bot/internal/api"
	"fluently/telegram-bot/internal/bot/handlers"
	"fluently/telegram-bot/internal/tasks"
	"fluently/telegram-bot/internal/webhook"
	"fluently/telegram-bot/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func main() {
	// Initialize configuration
	config.Init()
	cfg := config.GetConfig()

	// Initialize logger
	logger.Init(true) // true for development mode
	defer logger.Log.Sync()

	logger.Log.Info("Starting Fluently Telegram Bot",
		zap.String("version", "1.0.0"),
		zap.String("webhook_url", cfg.Bot.WebhookURL))

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	logger.Log.Info("Successfully connected to Redis", zap.String("addr", cfg.Redis.Addr))

	// Initialize API client
	apiClient := api.NewClient(cfg.API.BaseURL)

	// Initialize task scheduler
	scheduler := tasks.NewScheduler(cfg.Asynq.RedisAddr, cfg.Asynq.RedisPassword, cfg.Asynq.RedisDB)

	// Initialize Telegram bot (but don't start polling)
	bot, err := createTelegramBot(cfg)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	bot.Start()

	// Initialize handler service
	handlerService := handlers.NewHandlerService(cfg, redisClient, apiClient, scheduler, bot)

	// Initialize webhook server
	webhookServer := webhook.NewServer(cfg, handlerService)

	// Set up webhook with Telegram
	if err := setupWebhook(bot, cfg); err != nil {
		log.Fatalf("Failed to set up webhook: %v", err)
	}

	// Start background workers
	var wg sync.WaitGroup

	// Start Asynq worker server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startAsynqWorker(scheduler, handlerService); err != nil {
			logger.Log.Error("Asynq worker failed", zap.Error(err))
		}
	}()

	// Start webhook server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := webhookServer.Start(); err != nil {
			logger.Log.Error("Webhook server failed", zap.Error(err))
		}
	}()

	// Start cleanup routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		startCleanupRoutine(redisClient)
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Log.Info("Bot is running. Press Ctrl+C to stop.")
	<-sigChan

	logger.Log.Info("Shutting down...")

	// Close resources
	if err := scheduler.Close(); err != nil {
		logger.Log.Error("Failed to close scheduler", zap.Error(err))
	}

	if err := redisClient.Close(); err != nil {
		logger.Log.Error("Failed to close Redis client", zap.Error(err))
	}

	// Wait for goroutines to finish with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Log.Info("Graceful shutdown completed")
	case <-time.After(30 * time.Second):
		logger.Log.Warn("Shutdown timeout exceeded")
	}
}

// createTelegramBot creates and configures the Telegram bot
func createTelegramBot(cfg *config.Config) (*tele.Bot, error) {
	settings := tele.Settings{
		Token: cfg.Bot.Token,
	}

	bot, err := tele.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

// setupWebhook configures webhook with Telegram
func setupWebhook(bot *tele.Bot, cfg *config.Config) error {
	webhookConfig := &tele.Webhook{
		Listen: cfg.Webhook.Host + ":" + cfg.Webhook.Port,
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: cfg.Bot.WebhookURL,
		},
	}

	// Set webhook secret if configured
	if cfg.Webhook.Secret != "" {
		webhookConfig.SecretToken = cfg.Webhook.Secret
	}

	// Set TLS config if certificates are provided
	if cfg.Webhook.CertFile != "" && cfg.Webhook.KeyFile != "" {
		webhookConfig.TLS = &tele.WebhookTLS{
			Key:  cfg.Webhook.KeyFile,
			Cert: cfg.Webhook.CertFile,
		}
	}

	// Note: We're not calling bot.Start() with webhook because we're handling
	// webhook processing manually in our webhook server
	logger.Log.Info("Webhook configured",
		zap.String("url", cfg.Bot.WebhookURL),
		zap.String("listen", webhookConfig.Listen))

	return nil
}

// startAsynqWorker starts the Asynq worker server
func startAsynqWorker(scheduler *tasks.Scheduler, handlerService *handlers.HandlerService) error {
	server := scheduler.GetServer()

	// Register task handlers
	mux := setupTaskHandlers(handlerService)

	logger.Log.Info("Starting Asynq worker server")
	return server.Run(mux.GetServeMux())
}

// setupTaskHandlers sets up task handlers for Asynq
func setupTaskHandlers(handlerService *handlers.HandlerService) *tasks.TaskMux {
	mux := tasks.NewTaskMux()

	// Register task handlers
	mux.HandleFunc(tasks.TaskSendLessonReminder, handlerService.HandleLessonReminderTask)
	mux.HandleFunc(tasks.TaskSendDailyNotification, handlerService.HandleDailyNotificationTask)
	mux.HandleFunc(tasks.TaskGenerateLesson, handlerService.HandleGenerateLessonTask)
	mux.HandleFunc(tasks.TaskSyncProgress, handlerService.HandleSyncProgressTask)
	mux.HandleFunc(tasks.TaskCleanupSessions, handlerService.HandleCleanupSessionsTask)

	return mux
}

// startCleanupRoutine starts a routine to clean up expired sessions and data
func startCleanupRoutine(redisClient *redis.Client) {
	ticker := time.NewTicker(1 * time.Hour) // Run every hour
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

		// Clean up expired user sessions
		if err := cleanupExpiredSessions(ctx, redisClient); err != nil {
			logger.Log.Error("Failed to cleanup expired sessions", zap.Error(err))
		}

		// Clean up old progress data (older than 30 days)
		if err := cleanupOldProgressData(ctx, redisClient); err != nil {
			logger.Log.Error("Failed to cleanup old progress data", zap.Error(err))
		}

		cancel()
	}
}

// cleanupExpiredSessions removes expired user sessions from Redis
func cleanupExpiredSessions(ctx context.Context, redisClient *redis.Client) error {
	// Implementation would scan for session keys and remove expired ones
	logger.Log.Debug("Cleaning up expired sessions")

	// Example: scan for session:* keys and check expiration
	keys, err := redisClient.Keys(ctx, "session:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		ttl := redisClient.TTL(ctx, key).Val()
		if ttl < 0 { // Key has no expiration or expired
			redisClient.Del(ctx, key)
		}
	}

	logger.Log.Debug("Session cleanup completed", zap.Int("keys_checked", len(keys)))
	return nil
}

// cleanupOldProgressData removes old progress data to save memory
func cleanupOldProgressData(ctx context.Context, redisClient *redis.Client) error {
	logger.Log.Debug("Cleaning up old progress data")

	// Implementation would remove progress data older than threshold
	cutoff := time.Now().AddDate(0, 0, -30) // 30 days ago

	keys, err := redisClient.Keys(ctx, "user_progress:*").Result()
	if err != nil {
		return err
	}

	cleanedCount := 0
	for _, key := range keys {
		// Get the last activity time and check if it's too old
		// This would require parsing the progress data
		// For now, just check TTL
		ttl := redisClient.TTL(ctx, key).Val()
		if ttl < time.Duration(24*30)*time.Hour { // Less than 30 days TTL
			redisClient.Del(ctx, key)
			cleanedCount++
		}
	}

	logger.Log.Debug("Progress data cleanup completed",
		zap.Int("keys_checked", len(keys)),
		zap.Int("cleaned", cleanedCount),
		zap.Time("cutoff", cutoff))
	return nil
}
