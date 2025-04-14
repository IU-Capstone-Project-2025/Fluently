package postgres

import (
	"context"
	"errors"

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
}

func (r *WordPostgres) GetByID(ctx context.Context, idStr string) (*models.Word, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid uuid")
	}

	//TODO parameter query
	var word models.Word
	if err := r.db.WithContext(ctx).First(&word, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *WordPostgres) List(ctx context.Context) ([]*models.Word, error) {
	var words []*models.Word
	err := r.db.WithContext(ctx).Find(&words).Error

	return words, err
}

func (r *WordPostgres) Update(ctx context.Context, idStr string, updates map[string]any) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("invalid uuid")
	}

	return r.db.WithContext(ctx).Model(&models.Word{}).Where("id = ?", id).Updates(updates).Error
}

func (r *WordPostgres) Delete(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("invalid uuid")
	}

	return r.db.WithContext(ctx).Delete(&models.Word{}, "id = ?", id).Error
}
