package postgres

import (
	"context"
	"errors"
	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotLearnedWordsRepository is a repository for not learned words
type NotLearnedWordRepository struct {
	db *gorm.DB
}

// NewNotLearnedWordRepository creates a new instance of NotLearnedWordRepository
func NewNotLearnedWordRepository(db *gorm.DB) *NotLearnedWordRepository {
	return &NotLearnedWordRepository{db: db}
}

// Exists checks if a word is not learned
func (r *NotLearnedWordRepository) Exists(ctx context.Context, userID, wordID uuid.UUID) (bool, error) {
	var nlw models.NotLearnedWords
	err := r.db.WithContext(ctx).First(&nlw, "user_id = ? AND word_id = ?", userID, wordID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Create creates a new not learned word
func (r *NotLearnedWordRepository) Create(ctx context.Context, nlw *models.NotLearnedWords) error {
	return r.db.WithContext(ctx).Create(nlw).Error
}

func (r *NotLearnedWordRepository) DeleteIfExists(ctx context.Context, userID, wordID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND word_id = ?", userID, wordID).
		Delete(&models.NotLearnedWords{}).Error
}

// GetRecentlyNotLearnedWords returns recently not learned words for a user with word details
func (r *NotLearnedWordRepository) GetRecentlyNotLearnedWords(ctx context.Context, userID uuid.UUID, limit int) ([]models.Word, error) {
	var words []models.Word
	err := r.db.WithContext(ctx).
		Table("words").
		Joins("JOIN not_learned_words ON words.id = not_learned_words.word_id").
		Where("not_learned_words.user_id = ?", userID).
		Order("RANDOM()").
		Limit(limit).
		Find(&words).Error

	return words, err
}
