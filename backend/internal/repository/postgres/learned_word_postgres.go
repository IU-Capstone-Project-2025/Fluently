package postgres

import (
	"context"

<<<<<<< HEAD
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

func (r *LearnedWordRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.LearnedWords, error) {
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
	return r.db.Create(lw).Error
}

func (r *LearnedWordRepository) Update(ctx context.Context, lw *models.LearnedWords) error {
	return r.db.WithContext(ctx).Save(lw).Error
}

func (r *LearnedWordRepository) Delete(ctx context.Context, userID, wordID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&models.LearnedWords{}, "user_id = ? AND word_id = ?", userID, wordID).Error
=======
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fluently/go-backend/internal/repository/models"
)

type WordPostgres struct {
	db *gorm.DB
}

func NewWordPostgres(db *gorm.DB) *WordPostgres {
	return &WordPostgres{db: db}
}

func (r *WordPostgres) Create(ctx context.Context, word *models.Word) error {
	word.ID = uuid.New()
	return r.db.WithContext(ctx).Create(word).Error
>>>>>>> 514fbe1 (Add word create logic)
}
