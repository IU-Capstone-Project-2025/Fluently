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

func TestCreateUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	req := map[string]interface{}{
		"name":          "John Doe",
		"email":         "john@example.com",
		"provider":      "local",
		"google_id":     "google123",
		"password_hash": "secret",
		"role":          "user",
		"is_active":     true,
	}

	resp := e.POST("/users").
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

func TestGetUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := models.User{
		ID:           uuid.New(),
		Name:         "Jane",
		Email:        "jane@example.com",
		Role:         "admin",
		IsActive:     true,
		Provider:     "local",
		GoogleID:     "google456",
		PasswordHash: "hashed",
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	resp := e.GET("/users/" + user.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, "Jane", resp.Value("name").String().Raw())
	assert.Equal(t, "admin", resp.Value("role").String().Raw())
}

func TestUpdateUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := models.User{
		ID:           uuid.New(),
		Name:         "Initial",
		Email:        "init@example.com",
		Role:         "user",
		IsActive:     false,
		Provider:     "local",
		GoogleID:     "google789",
		PasswordHash: "hashed",
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"name":          "Updated Name",
		"email":         "updated@example.com",
		"provider":      "google",
		"google_id":     "google999",
		"password_hash": "newhash",
		"role":          "admin",
		"is_active":     true,
	}

	resp := e.PUT("/users/" + user.ID.String()).
		WithJSON(updateBody).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, "Updated Name", resp.Value("name").String().Raw())
	assert.Equal(t, "admin", resp.Value("role").String().Raw())
	assert.Equal(t, "updated@example.com", resp.Value("email").String().Raw())
}

func TestDeleteUser(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	user := models.User{
		ID:           uuid.New(),
		Name:         "DeleteMe",
		Email:        "deleteme@example.com",
		Role:         "user",
		IsActive:     true,
		Provider:     "local",
		GoogleID:     "gid",
		PasswordHash: "pass",
	}
	err := userRepo.Create(context.Background(), &user)
	assert.NoError(t, err)

	e.DELETE("/users/" + user.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	e.GET("/users/" + user.ID.String()).
		Expect().
		Status(http.StatusNotFound)
}
