package postgres

import (
	"context"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

func (r *PreferenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Preference, error) {
	var pref models.Preference
	if err := r.db.WithContext(ctx).First(&pref, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &pref, nil
}

func (r *PreferenceRepository) Update(ctx context.Context, pref *models.Preference) error {
	return r.db.WithContext(ctx).Save(pref).Error
}
