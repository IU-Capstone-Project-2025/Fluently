package postgres

import (
	"context"
	"errors"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LearnedWordRepository is a repository for learned words
type LearnedWordRepository struct {
	db *gorm.DB
}

// NewLearnedWordRepository creates a new instance of LearnedWordRepository
func NewLearnedWordRepository(db *gorm.DB) *LearnedWordRepository {
	return &LearnedWordRepository{db: db}
}

// ListByUserID returns a list of learned words for a user
func (r *LearnedWordRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.LearnedWords, error) {
	var words []models.LearnedWords
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&words).Error

	return words, err
}

// GetByUserWordID returns a learned word for a user
func (r *LearnedWordRepository) GetByUserWordID(ctx context.Context, userID, wordID uuid.UUID) (*models.LearnedWords, error) {
	var lw models.LearnedWords
	err := r.db.WithContext(ctx).
		First(&lw, "user_id = ? AND word_id = ?", userID, wordID).Error // learned word for user using word id and user id
	if err != nil {
		return nil, err
	}

	return &lw, nil
}

// IsLearned checks if a word is learned
func (r *LearnedWordRepository) IsLearned(ctx context.Context, userID, wordID uuid.UUID) (bool, error) {
	var lw models.LearnedWords
	err := r.db.WithContext(ctx).First(&lw, "user_id = ? AND word_id = ?", userID, wordID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // If the record is not found return false (not error)
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Create creates a new learned word
func (r *LearnedWordRepository) Create(ctx context.Context, lw *models.LearnedWords) error {
	return r.db.WithContext(ctx).Create(lw).Error
}

// Update updates a learned word
func (r *LearnedWordRepository) Update(ctx context.Context, lw *models.LearnedWords) error {
	return r.db.WithContext(ctx).Save(lw).Error
}

// Delete deletes a learned word
func (r *LearnedWordRepository) Delete(ctx context.Context, userID, wordID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&models.LearnedWords{}, "user_id = ? AND word_id = ?", userID, wordID).Error
}

// GetRecentlyLearnedWords returns recently learned words for a user with word details
func (r *LearnedWordRepository) GetRecentlyLearnedWords(ctx context.Context, userID uuid.UUID, limit int) ([]models.Word, error) {
	var words []models.Word
	err := r.db.WithContext(ctx).
		Table("words").
		Joins("JOIN learned_words ON words.id = learned_words.word_id").
		Where("learned_words.user_id = ?", userID).
		Order("learned_words.learned_at DESC").
		Limit(limit).
		Find(&words).Error

	return words, err
}
