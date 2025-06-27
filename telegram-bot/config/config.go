package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Bot    BotConfig
	Logger LoggerConfig
	Redis  RedisConfig
}

type BotConfig struct {
	Token string
}

type LoggerConfig struct {
	Level string
	Path  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var cfg *Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	viper.AutomaticEnv()

	cfg = &Config{
		Bot: BotConfig{
			Token: viper.GetString("BOT_TOKEN"),
		},
		Logger: LoggerConfig{
			Level: viper.GetString("LOG_LEVEL"),
			Path:  viper.GetString("LOG_PATH"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	}
}

func GetConfig() *Config {
	if cfg == nil {
		Init()
	}
	return cfg
}
