package main

import (
	"net/http"

	// Import docs only if they exist (conditional import for swag generation)
	"fluently/go-backend/docs"

	appConfig "fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/router"
	"fluently/go-backend/pkg/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Fluently API
// @version         1.0
// @description     Backend API for Fluently. Note: Auth routes are available at root level (/auth/*), while other API routes are under /api/v1/*
// @termsOfService  http://fluently-app.ru/terms/

// @contact.name   Danila Kochegarov
// @contact.url    http://fluently-app.ru
// @contact.email  Woolfer0097@yandex.ru

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host localhost:8070
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Config init
	appConfig.Init()
	logger.Init(true) // or false for production
	defer logger.Log.Sync()

	// Update Swagger host from configuration
	if swaggerHost := appConfig.GetConfig().Swagger.Host; swaggerHost != "" {
		docs.SwaggerInfo.Host = swaggerHost
		logger.Log.Info("Swagger host set from config", zap.String("host", swaggerHost))
	}

	// Database init
	dsn := appConfig.GetPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate
	err = db.AutoMigrate(
		&models.User{},
		&models.Preference{},
		&models.Topic{},
		&models.Word{},
		&models.Sentence{},
		&models.PickOption{},
		&models.RefreshToken{},
		&models.LearnedWords{},
		&models.LinkToken{},
	)
	if err != nil {
		logger.Log.Fatal("Failed to auto-migrate", zap.Error(err))
	}
	logger.Log.Info("Database migration completed successfully")

	//Init Router with routes
	r := chi.NewRouter()
	router.InitRoutes(db, r)

	// Start server
	logger.Log.Info("Logger initialization successful!")
	logger.Log.Info("App starting",
		zap.String("name", appConfig.GetConfig().API.AppName),
		zap.String("address", appConfig.GetConfig().API.AppHost+":"+appConfig.GetConfig().API.AppPort),
		zap.String("dsn", appConfig.GetPostgresDSN()),
	)

	err = http.ListenAndServe(appConfig.GetConfig().API.AppHost+":"+appConfig.GetConfig().API.AppPort, r)
	if err != nil {
		logger.Log.Fatal("App failed to start", zap.Error(err))
	}
}
