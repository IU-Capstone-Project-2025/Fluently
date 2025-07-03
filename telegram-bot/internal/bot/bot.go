package bot

import (
	"fluently/telegram-bot/config"
	"log"

	"github.com/redis/go-redis/v9"
	tele "gopkg.in/telebot.v3"
)

var RedisClient *redis.Client

func Start(cfg *config.Config, redisClient *redis.Client) {
	// Store Redis client globally for use in handlers
	RedisClient = redisClient

	pref := tele.Settings{
		Token: cfg.Bot.Token,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	bot.Start()
}
