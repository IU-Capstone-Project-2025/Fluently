package postgres

import (
	"context"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LearnedWordRepository struct {
	db *gorm.DB
}

func NewLearnedWordRepository(db *gorm.DB) *LearnedWordRepository {
	return &LearnedWordRepository{db: db}
}

func (r *LearnedWordRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.LearnedWords, error) {
	var words []models.LearnedWords
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&words).Error

	return words, err
}

func (r *LearnedWordRepository) GetByUserWordID(ctx context.Context, userID, wordID uuid.UUID) (*models.LearnedWords, error) {
	var lw models.LearnedWords
	err := r.db.WithContext(ctx).
		First(&lw, "user_id = ? AND word_id = ?", userID, wordID).Error
	if err != nil {
		return nil, err
	}

	return &lw, nil
}

func (r *LearnedWordRepository) Create(ctx context.Context, lw *models.LearnedWords) error {
	return r.db.WithContext(ctx).Create(lw).Error
}

func (r *LearnedWordRepository) Update(ctx context.Context, lw *models.LearnedWords) error {
	return r.db.WithContext(ctx).Save(lw).Error
}

func (r *LearnedWordRepository) Delete(ctx context.Context, userID, wordID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&models.LearnedWords{}, "user_id = ? AND word_id = ?", userID, wordID).Error
}
