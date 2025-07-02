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

	wordRepo *pg.WordRepository
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

	err = db.Exec("TRUNCATE TABLE words RESTART IDENTITY CASCADE").Error
	if err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}

	repo := pg.NewWordRepository(db)
	handler := &handlers.WordHandler{Repo: repo}

	r := chi.NewRouter()
	routes.RegisterWordRoutes(r, handler)

	testServer = httptest.NewServer(r)
}
