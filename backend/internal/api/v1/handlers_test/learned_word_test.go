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

// TestCreateLearnedWord tests the creation of a learned word
func TestCreateLearnedWord(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	user := models.User{
		ID:           uuid.New(),
		Email:        "test@user.com",
		Provider:     "local",
		PasswordHash: "hashed",
		Role:         "user",
		IsActive:     true,
	}
	assert.NoError(t, userRepo.Create(context.Background(), &user))

	// Set the test user for authentication
	setTestUser(&user)

	word := models.Word{
		ID:   uuid.New(),
		Word: "learn",
	}
	assert.NoError(t, wordRepo.Create(context.Background(), &word))

	req := map[string]interface{}{
		"user_id":          user.ID.String(),
		"word_id":          word.ID.String(),
		"learned_at":       time.Now().Format(time.RFC3339),
		"cnt_reviewed":     2,
		"confidence_score": 85,
	}

	e.POST("/api/v1/users/learned-words").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated)
}

// TestGetLearnedWord tests the retrieval of a learned word
func TestGetLearnedWord(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	user := models.User{ID: uuid.New(), Email: "get@test.com", Provider: "local", PasswordHash: "x", Role: "user", IsActive: true}
	word := models.Word{ID: uuid.New(), Word: "testword"}
	assert.NoError(t, userRepo.Create(context.Background(), &user))
	assert.NoError(t, wordRepo.Create(context.Background(), &word))

	// Set the test user for authentication
	setTestUser(&user)

	learned := models.LearnedWords{
		UserID:           user.ID,
		WordID:           word.ID,
		LearnedAt:        time.Now(),
		CountOfRevisions: 3,
		ConfidenceScore:  92,
	}
	assert.NoError(t, learnedWordRepo.Create(context.Background(), &learned))

	resp := e.GET("/api/v1/users/learned-words/" + word.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	resp.HasValue("user_id", user.ID.String())
	resp.HasValue("word_id", word.ID.String())
	resp.HasValue("cnt_reviewed", 3)
	resp.HasValue("confidence_score", 92)
	resp.Value("learned_at").String().NotEmpty()
}

// TestUpdateLearnedWord tests the update of a learned word
func TestUpdateLearnedWord(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	user := models.User{ID: uuid.New(), Email: "update@test.com", Provider: "local", PasswordHash: "x", Role: "user", IsActive: true}
	word := models.Word{ID: uuid.New(), Word: "updateword"}
	assert.NoError(t, userRepo.Create(context.Background(), &user))
	assert.NoError(t, wordRepo.Create(context.Background(), &word))

	// Set the test user for authentication
	setTestUser(&user)

	learned := models.LearnedWords{
		UserID:           user.ID,
		WordID:           word.ID,
		LearnedAt:        time.Now(),
		CountOfRevisions: 1,
		ConfidenceScore:  50,
	}
	assert.NoError(t, learnedWordRepo.Create(context.Background(), &learned))

	body := map[string]interface{}{
		"user_id":          user.ID.String(),
		"word_id":          word.ID.String(),
		"cnt_reviewed":     5,
		"confidence_score": 99,
	}
	e.PUT("/api/v1/users/learned-words/" + word.ID.String()).
		WithJSON(body).
		Expect().
		Status(http.StatusOK)

	updated, err := learnedWordRepo.GetByUserWordID(context.Background(), user.ID, word.ID)
	assert.NoError(t, err)
	assert.Equal(t, 5, updated.CountOfRevisions)
	assert.Equal(t, 99, updated.ConfidenceScore)
}

// TestDeleteLearnedWord tests the deletion of a learned word
func TestDeleteLearnedWord(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	user := models.User{ID: uuid.New(), Email: "delete@test.com", Provider: "local", PasswordHash: "x", Role: "user", IsActive: true}
	word := models.Word{ID: uuid.New(), Word: "deleteword"}
	assert.NoError(t, userRepo.Create(context.Background(), &user))
	assert.NoError(t, wordRepo.Create(context.Background(), &word))

	// Set the test user for authentication
	setTestUser(&user)

	learned := models.LearnedWords{
		UserID:           user.ID,
		WordID:           word.ID,
		LearnedAt:        time.Now(),
		CountOfRevisions: 1,
		ConfidenceScore:  60,
	}
	assert.NoError(t, learnedWordRepo.Create(context.Background(), &learned))

	e.DELETE("/api/v1/users/learned-words/" + word.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	_, err := learnedWordRepo.GetByUserWordID(context.Background(), user.ID, word.ID)
	assert.Error(t, err)
}
