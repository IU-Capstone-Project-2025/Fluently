package config

import (
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Auth     AuthConfig
	API      ApiConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	Google   GoogleConfig
	Swagger  SwaggerConfig
}

type AuthConfig struct {
	JWTSecret         string
	JWTExpiration     time.Duration
	RefreshExpiration time.Duration
	PasswordMinLength int
	RateLimitRequests int
	RateLimitDuration time.Duration
}

type ApiConfig struct {
	AppName string
	AppHost string
	AppPort string
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type LoggerConfig struct {
	Level string
	Path  string
}

type GoogleConfig struct {
	IosClientID     string
	AndroidClientID string
	WebClientID     string
}

type SwaggerConfig struct {
	AllowedEmails map[string]bool
}

var cfg *Config

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	// Tell Viper to automatically read env variables
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("APP_PORT", "8070")
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("REFRESH_EXPIRATION", "720h") // 30 days
	viper.SetDefault("PASSWORD_MIN_LENGTH", 8)
	viper.SetDefault("RATE_LIMIT_REQUESTS", 100)
	viper.SetDefault("RATE_LIMIT_DURATION", "1h")

	cfg = &Config{
		Auth: AuthConfig{
			JWTSecret:         viper.GetString("JWT_SECRET"),
			JWTExpiration:     viper.GetDuration("JWT_EXPIRATION"),
			RefreshExpiration: viper.GetDuration("REFRESH_EXPIRATION"),
			PasswordMinLength: viper.GetInt("PASSWORD_MIN_LENGTH"),
			RateLimitRequests: viper.GetInt("RATE_LIMIT_REQUESTS"),
			RateLimitDuration: viper.GetDuration("RATE_LIMIT_DURATION"),
		},
		API: ApiConfig{
			AppName: viper.GetString("APP_NAME"),
			AppHost: viper.GetString("APP_HOST"),
			AppPort: viper.GetString("APP_PORT"),
		},
		Database: DatabaseConfig{
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
		},
		Logger: LoggerConfig{
			Level: viper.GetString("LOG_LEVEL"),
			Path:  viper.GetString("LOG_PATH"),
		},
		Google: GoogleConfig{
			IosClientID:     viper.GetString("IOS_GOOGLE_CLIENT_ID"),
			AndroidClientID: viper.GetString("ANDROID_GOOGLE_CLIENT_ID"),
			WebClientID:     viper.GetString("WEB_GOOGLE_CLIENT_ID"),
		},
		Swagger: SwaggerConfig{
			AllowedEmails: parseEmailWhitelist(viper.GetString("SWAGGER_ALLOWED_EMAILS")),
		},
	}
}

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	if cfg == nil {
		Init()
	}
	return cfg
}

// GetPostgresDSN returns the PostgreSQL connection string
func GetPostgresDSN() string {
	return "postgres://" + cfg.Database.User + ":" + cfg.Database.Password +
		"@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Name
}

// GoogleOAuthConfig constructs an oauth2.Config based on application settings.
// This is used for mobile app OAuth flows, not for Swagger UI.
func GoogleOAuthConfig() *oauth2.Config {
	cfg := GetConfig()

	return &oauth2.Config{
		// RedirectURL will be set dynamically by the handler based on the client type
		ClientID:     cfg.Google.WebClientID,
		ClientSecret: os.Getenv("WEB_GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

// parseEmailWhitelist converts a comma-separated list into a map for O(1) lookups.
func parseEmailWhitelist(env string) map[string]bool {
	result := make(map[string]bool)
	for _, e := range strings.Split(env, ",") {
		e = strings.TrimSpace(e)
		if e != "" {
			result[e] = true
		}
	}
	return result
}
