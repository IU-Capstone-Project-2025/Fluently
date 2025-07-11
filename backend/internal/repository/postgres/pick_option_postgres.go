package postgres

import (
	"context"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PickOptionRepository struct {
	db *gorm.DB
}

func NewPickOptionRepository(db *gorm.DB) *PickOptionRepository {
	return &PickOptionRepository{db: db}
}

func (r *PickOptionRepository) Create(ctx context.Context, po *models.PickOption) error {
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *PickOptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PickOption, error) {
	var option models.PickOption
	if err := r.db.WithContext(ctx).First(&option, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &option, nil
}

func (r *PickOptionRepository) GetOptionByWordID(ctx context.Context, wordID uuid.UUID) (*models.PickOption, error) {
	var option models.PickOption
	if err := r.db.WithContext(ctx).First(&option, "word_id = ?", wordID).Error; err != nil {
		return nil, err
	}

	return &option, nil
}

func (r *PickOptionRepository) ListByWordID(ctx context.Context, wordID uuid.UUID) ([]models.PickOption, error) {
	var options []models.PickOption
	err := r.db.WithContext(ctx).
		Where("word_id = ?", wordID).
		Find(&options).Error

	return options, err
}

func (r *PickOptionRepository) Update(ctx context.Context, option *models.PickOption) error {
	return r.db.WithContext(ctx).Save(option).Error
}

func (r *PickOptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.PickOption{}, "id = ?", id).Error
}
