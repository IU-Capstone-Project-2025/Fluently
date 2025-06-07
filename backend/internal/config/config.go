package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	// Tell Viper to automatically read env variables
	viper.AutomaticEnv()

	// Optional: Set default values
	viper.SetDefault("APP_PORT", "3000")
}

// Getter helpers
func GetBotToken() string     { return viper.GetString("BOT_TOKEN") }
func GetAppName() string      { return viper.GetString("APP_NAME") }
func GetAppHost() string      { return viper.GetString("APP_HOST") }
func GetAppPort() string      { return viper.GetString("APP_PORT") }
func GetDBUser() string       { return viper.GetString("DB_USER") }
func GetDBPassword() string   { return viper.GetString("DB_PASSWORD") }
func GetDBHost() string       { return viper.GetString("DB_HOST") }
func GetDBPort() string       { return viper.GetString("DB_PORT") }
func GetDBName() string       { return viper.GetString("DB_NAME") }

// Optional: Construct DB DSN
func GetPostgresDSN() string {
	return "postgres://" + GetDBUser() + ":" + GetDBPassword() +
		"@" + GetDBHost() + ":" + GetDBPort() + "/" + GetDBName()
}
