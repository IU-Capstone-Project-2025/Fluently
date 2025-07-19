package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fluently/go-backend/internal/repository/models"
)

// TopicRepository is a repository for topics
type TopicRepository struct {
	db *gorm.DB
}

// NewTopicRepository creates a new instance of TopicRepository
func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

// GetByID returns a topic by id
func (r *TopicRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Topic, error) {
	var topic models.Topic
	if err := r.db.WithContext(ctx).First(&topic, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &topic, nil
}

// Create creates a new topic
func (r *TopicRepository) Create(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Create(topic).Error
}

// Update updates a topic
func (r *TopicRepository) Update(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Save(topic).Error
}

// Delete deletes a topic
func (r *TopicRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Topic{}, "id = ?", id).Error
}

// GetAll returns all topics
func (r *TopicRepository) GetAll(ctx context.Context) ([]models.Topic, error) {
	var topics []models.Topic
	if err := r.db.WithContext(ctx).Find(&topics).Error; err != nil {
		return nil, err
	}

	return topics, nil
}

// GetAllStartingWithCapital returns all topics that start with a capital letter
func (r *TopicRepository) GetAllStartingWithCapital(ctx context.Context) ([]models.Topic, error) {
	var topics []models.Topic
	if err := r.db.WithContext(ctx).
		Where("title ~ '^[A-Z]'").
		Find(&topics).Error; err != nil {
		return nil, err
	}

	return topics, nil
}
