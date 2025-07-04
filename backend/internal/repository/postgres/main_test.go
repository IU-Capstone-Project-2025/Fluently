package postgres

import (
	"os"
	"testing"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"

	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB

	userRepo        *UserRepository
	wordRepo        *WordRepository
	topicRepo       *TopicRepository
	sentenceRepo    *SentenceRepository
	preferenceRepo  *PreferenceRepository
	pickOptionRepo  *PickOptionRepository
	learnedWordRepo *LearnedWordRepository
)

func TestMain(m *testing.M) {
	dsn := config.GetPostgresDSNForTest()

	var err error
	db, err = gorm.Open(pgDriver.Open(dsn), &gorm.Config{})
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

	userRepo = NewUserRepository(db)
	wordRepo = NewWordRepository(db)
	topicRepo = NewTopicRepository(db)
	sentenceRepo = NewSentenceRepository(db)
	preferenceRepo = NewPreferenceRepository(db)
	pickOptionRepo = NewPickOptionRepository(db)
	learnedWordRepo = NewLearnedWordRepository(db)

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
