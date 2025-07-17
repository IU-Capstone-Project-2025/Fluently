package postgres

import (
	"context"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"

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

func (r *PreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Preference, error) {
	var pref models.Preference
	if err := r.db.WithContext(ctx).First(&pref, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &pref, nil
}



func (r *PreferenceRepository) Update(ctx context.Context, id uuid.UUID, req *schemas.UpdatePreferenceRequest) error {
	updates := map[string]interface{}{}

	if req.CEFRLevel != nil {
		updates["cefr_level"] = *req.CEFRLevel
	}
	if req.FactEveryday != nil {
		updates["fact_everyday"] = *req.FactEveryday
	}
	if req.Notifications != nil {
		updates["notifications"] = *req.Notifications
	}
	if req.NotificationAt != nil {
		updates["notifications_at"] = *req.NotificationAt
	}
	if req.WordsPerDay != nil {
		updates["words_per_day"] = *req.WordsPerDay
	}
	if req.Goal != nil {
		updates["goal"] = *req.Goal
	}
	if req.Subscribed != nil {
		updates["subscribed"] = *req.Subscribed
	}
	if req.AvatarImageURL != nil {
		updates["avatar_image_url"] = *req.AvatarImageURL
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Model(&models.Preference{}).
		Where("id = ?", id).
		Updates(updates).Error
}


func (r *PreferenceRepository) Create(ctx context.Context, pref *models.Preference) error {
	return r.db.WithContext(ctx).Create(pref).Error
}

func (r *PreferenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Preference{}, "id = ?", id).Error
}
