// main.go
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"telegram-bot/config"
	"telegram-bot/internal/api"
	"telegram-bot/internal/bot"
	"telegram-bot/internal/tasks"
)

func main() {
	// Parse command line flags
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Initialize configuration
	config.Init()
	cfg := config.GetConfig()

	// Initialize logger
	var logger *zap.Logger
	var err error
	if *debug || cfg.Logger.Level == "debug" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Starting Fluently Telegram Bot")

	// Initialize Redis client
	redisOpts := &redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	redisClient := redis.NewClient(redisOpts)

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Redis connection established")

	// Initialize API client
	apiClient := api.NewClient(cfg.API.BaseURL, logger)

	// Initialize task scheduler
	scheduler := tasks.NewScheduler(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, logger)

	// Create bot
	telegramBot, err := bot.NewTelegramBot(cfg, redisClient, apiClient, scheduler, logger)
	if err != nil {
		logger.Fatal("Failed to create telegram bot", zap.Error(err))
	}

	// Create a cancellable context
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the bot in a goroutine
	go func() {
		if err := telegramBot.Start(); err != nil {
			logger.Error("Bot stopped with error", zap.Error(err))
			cancel()
		}
	}()

	// Wait for termination signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	cancel()

	// Stop the bot
	telegramBot.Stop()
	logger.Info("Bot stopped gracefully")
}
