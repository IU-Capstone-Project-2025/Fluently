package main

import (
	"net/http"

	_ "fluently/go-backend/docs"
	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/router"
	"fluently/go-backend/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Fluently API
// @version         1.0
// @description     Backend API for Fluently Telegram bot
// @termsOfService  http://fluently.com/terms/

// @contact.name   Danila Kochegarov
// @contact.url    http://fluently.com
// @contact.email  Woolfer0097@yandex.ru

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	config.Init()
	logger.Init(true) // or false for production
	defer logger.Log.Sync()

	logger.Log.Info("Logger initialization successful!")
	logger.Log.Info("App starting",
		zap.String("name", config.GetAppName()),
		zap.String("address", config.GetAppHost()+":"+config.GetAppPort()),
		zap.String("dsn", config.GetPostgresDSN()),
	)

	db, err := gorm.Open(postgres.Open(config.GetPostgresDSN()), &gorm.Config{})

	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	r := router.InitRoutes(db)

	err = http.ListenAndServe(config.GetAppHost()+":"+config.GetAppPort(), r)
	if err != nil {
		logger.Log.Fatal("App failed to start", zap.Error(err))
	}
}
