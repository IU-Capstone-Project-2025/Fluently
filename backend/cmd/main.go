package main

import (
	"net/http"

	_ "fluently/go-backend/docs"
	"fluently/go-backend/internal/config"
	//"fluently/go-backend/internal/router"
	"fluently/go-backend/pkg/logger"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
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

	r := chi.NewRouter()

	// db, err := gorm.Open(postgres.Open(config.GetPostgresDSN()), &gorm.Config{})
	// router.InitRoutes(db)

	r.Get("/swagger/*", httpSwagger.WrapHandler) // Swagger UI

	r.Get("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	logger.Log.Info("Logger initialization successful!")
	logger.Log.Info("App starting",
		zap.String("name", config.GetAppName()),
		zap.String("address", config.GetAppHost()+":"+config.GetAppPort()),
		zap.String("dsn", config.GetPostgresDSN()),
	)

	err := http.ListenAndServe(config.GetAppHost()+":"+config.GetAppPort(), r)
	if err != nil {
		logger.Log.Fatal("App failed to start", zap.Error(err))
	}
}
