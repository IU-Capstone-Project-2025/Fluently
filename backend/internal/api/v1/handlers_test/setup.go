package handlers_test

import (
	"net/http/httptest"
	"testing"

	"fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/api/v1/routes"
	"fluently/go-backend/internal/repository/models"
	pg "fluently/go-backend/internal/repository/postgres"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	testServer *httptest.Server

	wordRepo     *pg.WordRepository
	userRepo     *pg.UserRepository
	topicRepo    *pg.TopicRepository
	sentenceRepo *pg.SentenceRepository
)

func setupTest(t *testing.T) {
	dsn := "host=localhost port=5433 user=test_user password=test_pass dbname=test_db sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	err = db.AutoMigrate(&models.Word{})
	if err != nil {
		t.Fatalf("failed to migrate DB: %v", err)
	}

	wordRepo = pg.NewWordRepository(db)
	userRepo = pg.NewUserRepository(db)
	topicRepo = pg.NewTopicRepository(db)
	sentenceRepo = pg.NewSentenceRepository(db)

	db.Exec("TRUNCATE TABLE words RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE topics RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE sentences RESTART IDENTITY CASCADE")

	wordHandler := &handlers.WordHandler{Repo: pg.NewWordRepository(db)}
	userHandler := &handlers.UserHandler{Repo: pg.NewUserRepository(db)}
	topicHandler := &handlers.TopicHandler{Repo: pg.NewTopicRepository(db)}
	sentenceHandler := &handlers.SentenceHandler{Repo: pg.NewSentenceRepository(db)}

	r := chi.NewRouter()
	routes.RegisterWordRoutes(r, wordHandler)
	routes.RegisterUserRoutes(r, userHandler)
	routes.RegisterTopicRoutes(r, topicHandler)
	routes.RegisterSentenceRoutes(r, sentenceHandler)

	testServer = httptest.NewServer(r)
}
