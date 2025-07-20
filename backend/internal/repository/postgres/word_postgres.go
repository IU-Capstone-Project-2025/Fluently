package postgres

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fluently/go-backend/internal/repository/models"
)

// WordRepository is a repository for words
type WordRepository struct {
	db *gorm.DB
}

// NewWordRepository creates a new instance of WordRepository
func NewWordRepository(db *gorm.DB) *WordRepository {
	return &WordRepository{db: db}
}

// ListWords returns all words
func (r *WordRepository) ListWords(ctx context.Context) ([]models.Word, error) {
	var words []models.Word
	if err := r.db.WithContext(ctx).Find(&words).Error; err != nil {
		return nil, err
	}

	return words, nil
}

// GetByID returns a word by id
func (r *WordRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Word, error) {
	var word models.Word
	if err := r.db.WithContext(ctx).First(&word, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &word, nil
}

// GetByValue returns a word by value
func (r *WordRepository) GetByValue(ctx context.Context, value string) (*models.Word, error) {
	var word models.Word
	if err := r.db.WithContext(ctx).First(&word, "word = ?", value).Error; err != nil {
		return nil, err
	}

	return &word, nil
}

// GetByWordTranslationPair returns a word by word and translation
func (r *WordRepository) GetByWordTranslationPair(ctx context.Context, word, translation string) (*models.Word, error) {
	var wordModel models.Word
	if err := r.db.WithContext(ctx).First(&wordModel, "word = ? AND translation = ?", word, translation).Error; err != nil {
		return nil, err
	}

	return &wordModel, nil
}

// GetRandomWordsByCEFRLevel returns random words by cefr level
func (r *WordRepository) GetRandomWordsByCEFRLevel(ctx context.Context, cefrLevel string, limit int) ([]models.Word, error) {
	var words []models.Word

	cefrLevel = strings.ToLower(cefrLevel)

	err := r.db.WithContext(ctx).
		Where("cefr_level = ?", cefrLevel).
		Order("RANDOM()").
		Limit(limit).
		Find(&words).Error

	if err != nil {
		return nil, err
	}

	return words, nil
}

// GetDayWord returns a random word by cefr level
func (r *WordRepository) GetDayWord(ctx context.Context, cefrLevel string, userID uuid.UUID) (*models.Word, error) {
	var word models.Word

	cefrLevel = strings.ToLower(cefrLevel)

	err := r.db.WithContext(ctx).
		Where("cefr_level = ?", cefrLevel).
		Order("RANDOM()").
		First(&word).Error

	if err != nil {
		return nil, err
	}

	return &word, nil
}

// Create creates a new word
func (r *WordRepository) Create(ctx context.Context, word *models.Word) error {
	return r.db.WithContext(ctx).Create(word).Error
}

// Update updates a word
func (r *WordRepository) Update(ctx context.Context, word *models.Word) error {
	return r.db.WithContext(ctx).Save(word).Error
}

// Delete deletes a word
func (r *WordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Word{}, "id = ?", id).Error
}

// GetRandomWordsWithTopic returns random words with topic information for topic/subtopic extraction
func (r *WordRepository) GetRandomWordsWithTopic(ctx context.Context, limit int) ([]models.Word, error) {
	var words []models.Word
	err := r.db.WithContext(ctx).
		Preload("Topic").
		Order("RANDOM()").
		Limit(limit).
		Find(&words).Error

	if err != nil {
		return nil, err
	}

	return words, err
}
