package postgres_test

import (
	"context"
	"testing"
	"time"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetPreference(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	now := time.Now()

	pref := &models.Preference{
		UserID:          userID,
		CEFRLevel:       "A2",
		FactEveryday:    true,
		Notifications:   true,
		NotificationsAt: &now,
		WordsPerDay:     15,
		Goal:            "Some Topic",
		Subscribed:      false,
		AvatarImageURL:  "http://banan.png",
	}

	err := preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)

	found, err := preferenceRepo.GetByID(ctx, pref.ID)
	assert.NoError(t, err)
	assert.Equal(t, pref.CEFRLevel, found.CEFRLevel)
	assert.Equal(t, pref.Goal, found.Goal)
	assert.Equal(t, pref.UserID, found.UserID)
	assert.Equal(t, pref.Subscribed, found.Subscribed)
}

func TestUpdatePreference(t *testing.T) {
	ctx := context.Background()

	userId := uuid.New()
	pref := &models.Preference{
		UserID:      userId,
		CEFRLevel:   "B1",
		WordsPerDay: 10,
		Goal:        "Some User Topic",
	}

	err := preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)

	pref.Goal = "Updated User Topic"
	pref.WordsPerDay = 20
	err = preferenceRepo.Update(ctx, pref)
	assert.NoError(t, err)

	updated, err := preferenceRepo.GetByID(ctx, pref.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User Topic", updated.Goal)
	assert.Equal(t, 20, updated.WordsPerDay)
}

func TestGetByUserID(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	pref := &models.Preference{
		UserID:      userID,
		CEFRLevel:   "C1",
		WordsPerDay: 7,
	}

	err := preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)

	found, err := preferenceRepo.GetByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, pref.ID, found.ID)
	assert.Equal(t, userID, found.UserID)
}

func TestDeletePreference(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	pref := &models.Preference{
		UserID:      userID,
		CEFRLevel:   "C1",
		WordsPerDay: 7,
	}

	err := preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)

	err = sentenceRepo.Delete(ctx, pref.ID)
	assert.NoError(t, err)

	_, err = sentenceRepo.GetByID(ctx, pref.ID)
	assert.Error(t, err) // "record not found"
}
