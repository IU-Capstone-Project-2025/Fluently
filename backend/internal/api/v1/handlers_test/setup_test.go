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

	// Drop all tables first to ensure clean state
	db.Exec("DROP TABLE IF EXISTS learned_words CASCADE")
	db.Exec("DROP TABLE IF EXISTS pick_options CASCADE")
	db.Exec("DROP TABLE IF EXISTS user_preferences CASCADE")
	db.Exec("DROP TABLE IF EXISTS sentences CASCADE")
	db.Exec("DROP TABLE IF EXISTS words CASCADE")
	db.Exec("DROP TABLE IF EXISTS topics CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Drop custom types if they exist to prevent conflicts
	db.Exec("DROP TYPE IF EXISTS string_array CASCADE")

	// Auto-migrate (UUID extension should already be available from init.sql)
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

	// Create handlers
	wordHandler := &handlers.WordHandler{Repo: pg.NewWordRepository(db)}
	userHandler := &handlers.UserHandler{Repo: pg.NewUserRepository(db)}
	topicHandler := &handlers.TopicHandler{Repo: pg.NewTopicRepository(db)}
	sentenceHandler := &handlers.SentenceHandler{Repo: pg.NewSentenceRepository(db)}
	prefHandler := &handlers.PreferenceHandler{Repo: pg.NewPreferenceRepository(db)}
	pickOptionHandler := &handlers.PickOptionHandler{Repo: pg.NewPickOptionRepository(db)}
	learnedWordHandler := &handlers.LearnedWordHandler{Repo: learnedWordRepo}
	progressHandler := &handlers.ProgressHandler{
		WordRepo:           wordRepo,
		LearnedWordRepo:    learnedWordRepo,
		NotLearnedWordRepo: pg.NewNotLearnedWordRepository(db),
		LLMClient:          utils.NewLLMClient(utils.LLMClientConfig{}),
		Redis:              utils.Redis(),
	}

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
