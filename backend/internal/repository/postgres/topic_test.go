package postgres

import (
	"context"
	"testing"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetTopic(t *testing.T) {
	ctx := context.Background()
	topic := &models.Topic{
		ID:    uuid.New(),
		Title: "Main topic",
	}

	err := topicRepo.Create(ctx, topic)
	assert.NoError(t, err)

	found, err := topicRepo.GetByID(ctx, topic.ID)
	assert.NoError(t, err)
	assert.Equal(t, topic.Title, found.Title)
}

func TestNestedTopics(t *testing.T) {
	ctx := context.Background()
	main := &models.Topic{
		ID:    uuid.New(),
		Title: "Root",
	}

	err := topicRepo.Create(ctx, main)
	assert.NoError(t, err)

	child := &models.Topic{
		ID:       uuid.New(),
		Title:    "Child",
		ParentID: &main.ID,
	}

	err = topicRepo.Create(ctx, child)
	assert.NoError(t, err)

	found, err := topicRepo.GetByID(ctx, child.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Child", found.Title)
	assert.NotNil(t, found.ParentID)
	assert.Equal(t, main.ID, *found.ParentID)
}
