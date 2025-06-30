package postgres_test

import (
	"os"
	"testing"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"

	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB

	userRepo        *postgres.UserRepository
	wordRepo        *postgres.WordRepository
	topicRepo       *postgres.TopicRepository
	sentenceRepo    *postgres.SentenceRepository
	preferenceRepo  *postgres.PreferenceRepository
	pickOptionRepo  *postgres.PickOptionRepository
	learnedWordRepo *postgres.LearnedWordRepository
)

func TestMain(m *testing.M) {
	dsn := "host=localhost port=5433 user=test_user password=test_pass dbname=test_db sslmode=disable"
	db, err := gorm.Open(pgDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
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
		panic("failed to migrate test database")
	}

	userRepo = postgres.NewUserRepository(db)
	wordRepo = postgres.NewWordRepository(db)
	topicRepo = postgres.NewTopicRepository(db)
	sentenceRepo = postgres.NewSentenceRepository(db)
	preferenceRepo = postgres.NewPreferenceRepository(db)
	pickOptionRepo = postgres.NewPickOptionRepository(db)
	learnedWordRepo = postgres.NewLearnedWordRepository(db)

	// Clear all tables before test
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE words RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE topics RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE sentences RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE user_preferences RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE pick_options RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE learned_words RESTART IDENTITY CASCADE")

	code := m.Run()

	os.Exit(code)
}
