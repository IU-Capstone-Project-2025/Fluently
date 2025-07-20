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

// TestCreateUser tests the creation of a new user
func TestCreateUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	req := map[string]interface{}{
		"name":          "John Doe",
		"email":         "john-" + uuid.New().String()[:8] + "@example.com",
		"provider":      "local",
		"google_id":     "google123",
		"password_hash": "secret",
		"role":          "user",
		"is_active":     true,
	}

	resp := e.POST("/api/v1/users").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.Equal(t, "John Doe", resp.Value("name").String().Raw())
	assert.Equal(t, "john@example.com", resp.Value("email").String().Raw())
	assert.Equal(t, "user", resp.Value("role").String().Raw())
	assert.Equal(t, true, resp.Value("is_active").Boolean().Raw())
	assert.NotEmpty(t, resp.Value("id").String().Raw())
	assert.NotEmpty(t, resp.Value("created_at").String().Raw())
}

// TestGetUser tests the retrieval of a user
func TestGetUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := models.User{
		ID:           uuid.New(),
		Name:         "Jane",
		Email:        "jane-" + uuid.New().String()[:8] + "@example.com",
		Role:         "admin",
		IsActive:     true,
		Provider:     "local",
		GoogleID:     "google456",
		PasswordHash: "hashed",
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/users/" + user.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, "Jane", resp.Value("name").String().Raw())
	assert.Equal(t, "admin", resp.Value("role").String().Raw())
}

// TestUpdateUser tests the update of a user
func TestUpdateUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := models.User{
		ID:           uuid.New(),
		Name:         "Original",
		Email:        "original-" + uuid.New().String()[:8] + "@example.com",
		Role:         "user",
		IsActive:     true,
		Provider:     "local",
		GoogleID:     "google789",
		PasswordHash: "hashed",
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"name":      "Updated Name",
		"email":     "updated-" + uuid.New().String()[:8] + "@example.com",
		"role":      "admin",
		"is_active": true,
	}

	resp := e.PUT("/api/v1/users/" + user.ID.String()).
		WithJSON(updateBody).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, "Updated Name", resp.Value("name").String().Raw())
	assert.Equal(t, "admin", resp.Value("role").String().Raw())
	assert.Equal(t, "updated@example.com", resp.Value("email").String().Raw())
}

// TestDeleteUser tests the deletion of a user
func TestDeleteUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := models.User{
		ID:           uuid.New(),
		Name:         "Delete Me",
		Email:        "delete-" + uuid.New().String()[:8] + "@example.com",
		Role:         "user",
		IsActive:     true,
		Provider:     "local",
		GoogleID:     "google999",
		PasswordHash: "hashed",
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	e.DELETE("/api/v1/users/" + user.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	e.GET("/api/v1/users/" + user.ID.String()).
		Expect().
		Status(http.StatusNotFound)
}
