package postgres

import (
	"context"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SentenceRepository is a repository for sentences
type SentenceRepository struct {
	db *gorm.DB
}

// NewSentenceRepository creates a new instance of SentenceRepository
func NewSentenceRepository(db *gorm.DB) *SentenceRepository {
	return &SentenceRepository{db: db}
}

// ListByWord returns a list of sentences for a word
func (r *SentenceRepository) ListByWord(ctx context.Context, wordID uuid.UUID) ([]models.Sentence, error) {
	var sentences []models.Sentence
	err := r.db.WithContext(ctx).Where("word_id = ?", wordID).Find(&sentences).Error

	return sentences, err
}

// Create creates a new sentence
func (r *SentenceRepository) Create(ctx context.Context, s *models.Sentence) error {
	return r.db.WithContext(ctx).Create(s).Error
}

// Update updates a sentence
func (r *SentenceRepository) Update(ctx context.Context, s *models.Sentence) error {
	return r.db.WithContext(ctx).Save(s).Error
}

// Delete deletes a sentence
func (r *SentenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Sentence{}, "id = ?", id).Error
}

// GetByID returns a sentence by id
func (r *SentenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Sentence, error) {
	var s models.Sentence
	err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// GetByWordID returns a list of sentences for a word
func (r *SentenceRepository) GetByWordID(ctx context.Context, wordID uuid.UUID) ([]models.Sentence, error) {
	var sentences []models.Sentence
	err := r.db.WithContext(ctx).Find(&sentences, "word_id = ?", wordID).Error
	if err != nil {
		return nil, err
	}

	return sentences, nil
}
