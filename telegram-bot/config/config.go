package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Bot    BotConfig
	Logger LoggerConfig
	Redis  RedisConfig
	API    APIConfig
	Asynq  AsynqConfig
	TTS    TTSConfig
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

type APIConfig struct {
	BaseURL string
	APIKey  string
	Timeout int
}

type AsynqConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
}

type TTSConfig struct {
	CacheDir string
}

var cfg *Config

func Init() {
	// Try to find .env file in multiple locations
	envPaths := []string{
		".env",          // Current directory
		"../.env",       // Parent directory (if running from telegram-bot/)
		"../../.env",    // Two levels up (if running from telegram-bot/cmd/)
		"../../../.env", // Three levels up (if running from telegram-bot/config/)
	}

	var envLoaded bool
	for _, envPath := range envPaths {
		if absPath, err := filepath.Abs(envPath); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				if err := godotenv.Load(envPath); err == nil {
					log.Printf("Loaded .env file from: %s", absPath)
					envLoaded = true
					break
				}
			}
		}
	}

	if !envLoaded {
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

		API: APIConfig{
			BaseURL: viper.GetString("API_BASE_URL"),
			APIKey:  viper.GetString("API_KEY"),
			Timeout: viper.GetInt("API_TIMEOUT"),
		},
		Asynq: AsynqConfig{
			RedisAddr:     viper.GetString("ASYNQ_REDIS_ADDR"),
			RedisPassword: viper.GetString("ASYNQ_REDIS_PASSWORD"),
			RedisDB:       viper.GetInt("ASYNQ_REDIS_DB"),
			Concurrency:   viper.GetInt("ASYNQ_CONCURRENCY"),
		},
		TTS: TTSConfig{
			CacheDir: viper.GetString("TTS_CACHE_DIR"),
		},
	}

	// Set defaults
	if cfg.API.Timeout == 0 {
		cfg.API.Timeout = 30
	}
	if cfg.Asynq.Concurrency == 0 {
		cfg.Asynq.Concurrency = 10
	}
	if cfg.Asynq.RedisAddr == "" {
		cfg.Asynq.RedisAddr = cfg.Redis.Addr
	}
	if cfg.Asynq.RedisPassword == "" {
		cfg.Asynq.RedisPassword = cfg.Redis.Password
	}
	if cfg.Asynq.RedisDB == 0 {
		cfg.Asynq.RedisDB = cfg.Redis.DB
	}
	if cfg.TTS.CacheDir == "" {
		cfg.TTS.CacheDir = "/tmp/tts"
	}
}

func GetConfig() *Config {
	if cfg == nil {
		Init()
	}
	return cfg
}
