package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"fluently/go-backend/internal/repository/models"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCreateUserPreferences tests the creation of user preferences
func TestCreateUserPreferences(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	userID := uuid.New()

	// Create the user first to avoid foreign key constraint violation
	user := models.User{
		ID:           userID,
		Email:        "preferences@test.com",
		Provider:     "local",
		PasswordHash: "hashed",
		Role:         "user",
		IsActive:     true,
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	req := map[string]interface{}{
		"user_id":          userID.String(),
		"cefr_level":       "B2",
		"fact_everyday":    true,
		"notifications":    true,
		"words_per_day":    20,
		"goal":             "Improve vocabulary",
		"subscribed":       true,
		"avatar_image_url": "http://example.com/avatar.png",
	}

	resp := e.POST("/api/v1/preferences/user/" + userID.String()).
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.Equal(t, userID.String(), resp.Value("user_id").String().Raw())
	assert.Equal(t, "B2", resp.Value("cefr_level").String().Raw())
	assert.Equal(t, true, resp.Value("fact_everyday").Boolean().Raw())
	assert.Equal(t, true, resp.Value("notifications").Boolean().Raw())
	assert.Equal(t, float64(20), resp.Value("words_per_day").Number().Raw())
	assert.Equal(t, "Improve vocabulary", resp.Value("goal").String().Raw())
	assert.Equal(t, true, resp.Value("subscribed").Boolean().Raw())
	assert.Equal(t, "http://example.com/avatar.png", resp.Value("avatar_image_url").String().Raw())
}

// TestDeletePreference tests the deletion of user preferences
func TestDeletePreference(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	userID := uuid.New()

	// Create the user first to avoid foreign key constraint violation
	user := models.User{
		ID:           userID,
		Email:        "delete-pref@test.com",
		Provider:     "local",
		PasswordHash: "hashed",
		Role:         "user",
		IsActive:     true,
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	pref := models.Preference{
		ID:     uuid.New(),
		UserID: userID,
	}
	err = prefRepo.Create(context.Background(), &pref)
	assert.NoError(t, err)

	// Set the test user context for authentication
	setTestUser(&user)

	e.DELETE("/api/v1/preferences").
		Expect().
		Status(http.StatusNoContent)

	// Verify it was deleted by trying to get it by user ID
	_, err = prefRepo.GetByUserID(context.Background(), userID)
	assert.Error(t, err, "Expected error when trying to get deleted preference")
}
