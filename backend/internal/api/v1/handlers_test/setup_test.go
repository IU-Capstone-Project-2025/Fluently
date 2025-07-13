package handlers_test

import (
	"net/http/httptest"
	"testing"

	"fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/api/v1/routes"
	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"
	pg "fluently/go-backend/internal/repository/postgres"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	testServer *httptest.Server

	wordRepo        *pg.WordRepository
	userRepo        *pg.UserRepository
	topicRepo       *pg.TopicRepository
	sentenceRepo    *pg.SentenceRepository
	prefRepo        *pg.PreferenceRepository
	pickOptionRepo  *pg.PickOptionRepository
	learnedWordRepo *pg.LearnedWordRepository
)

func setupTest(t *testing.T) {
	dsn := config.GetPostgresDSNForTest()
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	err = db.AutoMigrate(
		&models.User{},
		&models.Word{},
		&models.Topic{},
		&models.Sentence{},
		&models.Preference{},
		&models.PickOption{},
		&models.LearnedWords{},
	)
	if err != nil {
		t.Fatalf("failed to migrate DB: %v", err)
	}

	wordRepo = pg.NewWordRepository(db)
	userRepo = pg.NewUserRepository(db)
	topicRepo = pg.NewTopicRepository(db)
	sentenceRepo = pg.NewSentenceRepository(db)
	prefRepo = pg.NewPreferenceRepository(db)
	pickOptionRepo = pg.NewPickOptionRepository(db)
	learnedWordRepo = pg.NewLearnedWordRepository(db)

	// Clear all tables in proper dependency order (dependent tables first)
	db.Exec("TRUNCATE TABLE learned_words RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE pick_options RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE user_preferences RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE sentences RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE words RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE topics RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	wordHandler := &handlers.WordHandler{Repo: pg.NewWordRepository(db)}
	userHandler := &handlers.UserHandler{Repo: pg.NewUserRepository(db)}
	topicHandler := &handlers.TopicHandler{Repo: pg.NewTopicRepository(db)}
	sentenceHandler := &handlers.SentenceHandler{Repo: pg.NewSentenceRepository(db)}
	prefHandler := &handlers.PreferenceHandler{Repo: pg.NewPreferenceRepository(db)}
	pickOptionHandler := &handlers.PickOptionHandler{Repo: pg.NewPickOptionRepository(db)}
	learnedWordHandler := &handlers.LearnedWordHandler{Repo: learnedWordRepo}
	progressHandler := &handlers.ProgressHandler{WordRepo: wordRepo, LearnedWordRepo: learnedWordRepo}

	r := chi.NewRouter()
	routes.RegisterWordRoutes(r, wordHandler)
	routes.RegisterUserRoutes(r, userHandler)
	routes.RegisterTopicRoutes(r, topicHandler)
	routes.RegisterSentenceRoutes(r, sentenceHandler)
	routes.RegisterPreferencesRoutes(r, prefHandler)
	routes.RegisterPickOptionRoutes(r, pickOptionHandler)
	routes.RegisterLearnedWordRoutes(r, learnedWordHandler)
	routes.RegisterProgressRoutes(r, progressHandler)

	testServer = httptest.NewServer(r)
}
