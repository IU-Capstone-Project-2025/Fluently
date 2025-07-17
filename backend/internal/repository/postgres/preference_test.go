package postgres

import (
	"context"
	"testing"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCreateAndGetPreference is a test for creating and getting a preference
func TestCreateAndGetPreference(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:       uuid.New(),
		Name:     "Pref User",
		Email:    "pref@example.com",
		Role:     "user",
		IsActive: true,
	}

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	now := time.Now()
	pref := &models.Preference{
		UserID:          user.ID,
		CEFRLevel:       "A2",
		FactEveryday:    true,
		Notifications:   true,
		NotificationsAt: &now,
		WordsPerDay:     15,
		Goal:            "Some Topic",
		Subscribed:      false,
		AvatarImageURL:  "http://banan.png",
	}

	err = preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)

	found, err := preferenceRepo.GetByID(ctx, pref.ID)
	assert.NoError(t, err)
	assert.Equal(t, pref.CEFRLevel, found.CEFRLevel)
	assert.Equal(t, pref.Goal, found.Goal)
	assert.Equal(t, pref.UserID, found.UserID)
	assert.Equal(t, pref.Subscribed, found.Subscribed)
}

// TestUpdatePreference is a test for updating a preference
func TestUpdatePreference(t *testing.T) {
	ctx := context.Background()
	var req schemas.UpdatePreferenceRequest

	user := &models.User{
		ID:       uuid.New(),
		Name:     "UpdatePref User",
		Email:    "updatepref@example.com",
		Role:     "user",
		IsActive: true,
	}

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	pref := &models.Preference{
		UserID:      user.ID,
		CEFRLevel:   "B1",
		WordsPerDay: 10,
		Goal:        "Some User Topic",
	}

	err = preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, pref.ID)

	updatedGoal := "Updated User Topic"
	updatedWordPerDay := 20

	req.Goal = &updatedGoal
	req.WordsPerDay = &updatedWordPerDay

	err = preferenceRepo.Update(ctx, pref.ID, &req)
	assert.NoError(t, err)

	updated, err := preferenceRepo.GetByID(ctx, pref.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User Topic", updated.Goal)
	assert.Equal(t, 20, updated.WordsPerDay)
}

// TestGetByUserID is a test for getting a preference by user
func TestGetByUserID(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:       uuid.New(),
		Name:     "GetPrefByUser User",
		Email:    "GetPrefByUser@example.com",
		Role:     "user",
		IsActive: true,
	}

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	pref := &models.Preference{
		UserID:      user.ID,
		CEFRLevel:   "C1",
		WordsPerDay: 7,
	}

	err = preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)

	found, err := preferenceRepo.GetByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, pref.ID, found.ID)
	assert.Equal(t, user.ID, found.UserID)
}

// TestDeletePreference is a test for deleting a preference
func TestDeletePreference(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:       uuid.New(),
		Name:     "DeletePref User",
		Email:    "DeletePref@example.com",
		Role:     "user",
		IsActive: true,
	}

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	pref := &models.Preference{
		UserID:      user.ID,
		CEFRLevel:   "C1",
		WordsPerDay: 7,
	}

	err = preferenceRepo.Create(ctx, pref)
	assert.NoError(t, err)

	err = preferenceRepo.Delete(ctx, pref.ID)
	assert.NoError(t, err)

	_, err = preferenceRepo.GetByID(ctx, pref.ID)
	assert.Error(t, err) // "record not found"
}
