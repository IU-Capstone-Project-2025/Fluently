package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fluently/go-backend/internal/repository/models"
)

type WordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) ListWords(ctx context.Context) ([]models.Word, error) {
	var words []models.Word
	if err := r.db.WithContext(ctx).Find(&words).Error; err != nil {
		return nil, err
	}

	return words, nil
}

func (r *WordRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Word, error) {
	var word models.Word
	if err := r.db.WithContext(ctx).First(&word, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *WordRepository) GetByValue(ctx context.Context, value string) (*models.Word, error) {
	var word models.Word
	if err := r.db.WithContext(ctx).First(&word, "word = ?", value).Error; err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *WordRepository) GetByWordTranslationPair(ctx context.Context, word, translation string) (*models.Word, error) {
	var wordModel models.Word
	if err := r.db.WithContext(ctx).First(&wordModel, "word = ? AND translation = ?", word, translation).Error; err != nil {
		return nil, err
	}

	return &wordModel, nil
}

func (r *WordRepository) Create(ctx context.Context, word *models.Word) error {
	return r.db.WithContext(ctx).Create(word).Error
}

func (r *WordRepository) Update(ctx context.Context, word *models.Word) error {
	return r.db.WithContext(ctx).Save(word).Error
}

func (r *WordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Word{}, "id = ?", id).Error
}
