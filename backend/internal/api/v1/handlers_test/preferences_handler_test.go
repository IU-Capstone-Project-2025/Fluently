package handlers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"fluently/go-backend/internal/repository/models"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserPreferences(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := &models.User{
		ID:    uuid.New(),
		Email: "test1@example.com",
	}
	err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	reqBody := map[string]interface{}{
		"cefr_level":       "B2",
		"fact_everyday":    true,
		"notifications":    true,
		"notification_at":  time.Now().Format(time.RFC3339),
		"words_per_day":    20,
		"goal":             "Improve vocabulary",
		"subscribed":       true,
		"avatar_image_url": "http://example.com/avatar.png",
	}

	resp := e.POST("/preferences/" + user.ID.String() + "/").
		WithJSON(reqBody).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.Equal(t, user.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, user.ID.String(), resp.Value("user_id").String().Raw())
	assert.Equal(t, "B2", resp.Value("cefr_level").String().Raw())
	assert.Equal(t, true, resp.Value("fact_everyday").Raw())
	assert.Equal(t, true, resp.Value("notifications").Raw())
	assert.Equal(t, 20, int(resp.Value("words_per_day").Number().Raw()))
	assert.Equal(t, "Improve vocabulary", resp.Value("goal").String().Raw())
	assert.Equal(t, true, resp.Value("subscribed").Raw())
	assert.Equal(t, "http://example.com/avatar.png", resp.Value("avatar_image_url").String().Raw())
}

// func TestGetUserPreferences(t *testing.T) {
// 	setupTest(t)

// 	e := httpexpect.Default(t, testServer.URL)

// 	user := &models.User{
// 		ID:    uuid.New(),
// 		Email: "test2@example.com",
// 	}
// 	err := userRepo.Create(context.Background(), user)
// 	assert.NoError(t, err)

// 	pref := models.Preference{
// 		ID:              user.ID,
// 		UserID:          user.ID,
// 		CEFRLevel:       "C1",
// 		FactEveryday:    false,
// 		Notifications:   false,
// 		NotificationsAt: nil,
// 		WordsPerDay:     15,
// 		Goal:            "Learn daily",
// 		Subscribed:      false,
// 		AvatarImageURL:  "",
// 	}
// 	err = prefRepo.Create(context.Background(), &pref)
// 	assert.NoError(t, err)

// 	resp := e.GET("/preferences/" + user.ID.String() + "/").
// 		Expect().
// 		Status(http.StatusOK).
// 		JSON().Object()

// 	assert.Equal(t, user.ID.String(), resp.Value("id").String().Raw())
// 	assert.Equal(t, "C1", resp.Value("cefr_level").String().Raw())
// 	assert.Equal(t, false, resp.Value("fact_everyday").Raw())
// 	assert.Equal(t, false, resp.Value("notifications").Raw())
// 	assert.Equal(t, 15, int(resp.Value("words_per_day").Number().Raw()))
// 	assert.Equal(t, "Learn daily", resp.Value("goal").String().Raw())
// 	assert.Equal(t, false, resp.Value("subscribed").Raw())
// 	assert.Equal(t, "", resp.Value("avatar_image_url").String().Raw())
// }

func TestUpdateUserPreferences(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := &models.User{
		ID:    uuid.New(),
		Email: "test3@example.com",
	}
	err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	pref := models.Preference{
		ID:            user.ID,
		UserID:        user.ID,
		CEFRLevel:     "B1",
		FactEveryday:  false,
		Notifications: false,
		WordsPerDay:   10,
		Goal:          "Initial goal",
		Subscribed:    false,
	}
	err = prefRepo.Create(context.Background(), &pref)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"cefr_level":       "C2",
		"fact_everyday":    true,
		"notifications":    true,
		"notification_at":  time.Now().Format(time.RFC3339),
		"words_per_day":    25,
		"goal":             "Updated goal",
		"subscribed":       true,
		"avatar_image_url": "http://example.com/new_avatar.png",
	}

	resp := e.PUT("/preferences/" + user.ID.String() + "/").
		WithJSON(updateBody).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, user.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, "C2", resp.Value("cefr_level").String().Raw())
	assert.Equal(t, true, resp.Value("fact_everyday").Raw())
	assert.Equal(t, true, resp.Value("notifications").Raw())
	assert.Equal(t, 25, int(resp.Value("words_per_day").Number().Raw()))
	assert.Equal(t, "Updated goal", resp.Value("goal").String().Raw())
	assert.Equal(t, true, resp.Value("subscribed").Raw())
	assert.Equal(t, "http://example.com/new_avatar.png", resp.Value("avatar_image_url").String().Raw())
}

func TestDeletePreference(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := &models.User{
		ID:    uuid.New(),
		Email: "test4@example.com",
	}
	err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	pref := models.Preference{
		ID:     user.ID,
		UserID: user.ID,
	}
	err = prefRepo.Create(context.Background(), &pref)
	assert.NoError(t, err)

	e.DELETE("/preferences/" + user.ID.String() + "/").
		Expect().
		Status(http.StatusNoContent)

	_, err = prefRepo.GetByID(context.Background(), user.ID)
	assert.Error(t, err)
}
