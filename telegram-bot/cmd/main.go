// main.go
package main

import (
	"context"
	"log"

	"fluently/telegram-bot/config"
	"fluently/telegram-bot/internal/bot"
	"fluently/telegram-bot/pkg/logger"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize configuration
	config.Init()
	cfg := config.GetConfig()

	// Initialize logger
	logger.Init(true) // true for development mode

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
	log.Println("Successfully connected to Redis")

	// Start the bot
	bot.Start(cfg, redisClient)
}
