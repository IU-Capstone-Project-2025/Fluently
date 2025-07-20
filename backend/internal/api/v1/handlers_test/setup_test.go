package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/api/v1/routes"
	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/middleware"
	"fluently/go-backend/internal/repository/models"
	pg "fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTest sets up the test server
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

	// Global test user for authentication
	testUser      *models.User
	testUserMutex sync.RWMutex
)

// setupTest sets up the test server
// This function is called before each test
// It sets up the test server and the database
func setupTest(t *testing.T) {
	dsn := config.GetPostgresDSNForTest()
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	// Connect extension for UUID
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Auto-migrate
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

	// Create repositories
	wordRepo = pg.NewWordRepository(db)
	userRepo = pg.NewUserRepository(db)
	topicRepo = pg.NewTopicRepository(db)
	sentenceRepo = pg.NewSentenceRepository(db)
	prefRepo = pg.NewPreferenceRepository(db)
	pickOptionRepo = pg.NewPickOptionRepository(db)
	learnedWordRepo = pg.NewLearnedWordRepository(db)

	// Clear all tables in proper dependency order (dependent tables first)
	// Use DELETE instead of TRUNCATE to avoid deadlocks in concurrent tests
	db.Exec("DELETE FROM learned_words")
	db.Exec("DELETE FROM pick_options")
	db.Exec("DELETE FROM user_preferences")
	db.Exec("DELETE FROM sentences")
	db.Exec("DELETE FROM words")
	db.Exec("DELETE FROM topics")
	db.Exec("DELETE FROM users")

	// Reset sequences after DELETE operations
	db.Exec("ALTER SEQUENCE IF EXISTS learned_words_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE IF EXISTS pick_options_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE IF EXISTS user_preferences_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE IF EXISTS sentences_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE IF EXISTS words_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE IF EXISTS topics_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 1")

	// Create handlers
	wordHandler := &handlers.WordHandler{Repo: pg.NewWordRepository(db)}
	userHandler := &handlers.UserHandler{Repo: pg.NewUserRepository(db)}
	topicHandler := &handlers.TopicHandler{Repo: pg.NewTopicRepository(db)}
	sentenceHandler := &handlers.SentenceHandler{Repo: pg.NewSentenceRepository(db)}
	prefHandler := &handlers.PreferenceHandler{Repo: pg.NewPreferenceRepository(db)}
	pickOptionHandler := &handlers.PickOptionHandler{Repo: pg.NewPickOptionRepository(db)}
	learnedWordHandler := &handlers.LearnedWordHandler{Repo: learnedWordRepo}
	progressHandler := &handlers.ProgressHandler{WordRepo: wordRepo, LearnedWordRepo: learnedWordRepo}

	// Create router
	r := chi.NewRouter()

	// Initialize JWT auth for tests
	utils.InitJWTAuth()

	// Add authentication middleware for protected routes
	r.Route("/api/v1", func(r chi.Router) {
		// For tests, we'll use a simplified auth middleware that doesn't require real JWT
		r.Use(testAuthMiddleware)

		// Protected API routes
		routes.RegisterUserRoutes(r, userHandler)
		routes.RegisterWordRoutes(r, wordHandler)
		routes.RegisterTopicRoutes(r, topicHandler)
		routes.RegisterSentenceRoutes(r, sentenceHandler)
		routes.RegisterPreferencesRoutes(r, prefHandler)
		routes.RegisterPickOptionRoutes(r, pickOptionHandler)
		routes.RegisterLearnedWordRoutes(r, learnedWordHandler)
		routes.RegisterProgressRoutes(r, progressHandler)
	})

	// Create test server
	testServer = httptest.NewServer(r)
}

// setTestUser sets the user for the current test
func setTestUser(user *models.User) {
	testUserMutex.Lock()
	defer testUserMutex.Unlock()
	testUser = user
}

// testAuthMiddleware is a simplified authentication middleware for tests
func testAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testUserMutex.RLock()
		user := testUser
		testUserMutex.RUnlock()

		if user == nil {
			// Create a default test user if none is set
			user = &models.User{
				ID:    uuid.New(),
				Email: "test@example.com",
				Role:  "user",
			}
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), middleware.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
