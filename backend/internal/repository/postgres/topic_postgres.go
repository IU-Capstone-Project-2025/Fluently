package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fluently/go-backend/internal/repository/models"
)

type TopicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Topic, error) {
	var topic models.Topic
	if err := r.db.WithContext(ctx).First(&topic, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &topic, nil
}

func (r *TopicRepository) Create(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *TopicRepository) Update(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Save(topic).Error
}

func (r *TopicRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Topic{}, "id = ?", id).Error
}
