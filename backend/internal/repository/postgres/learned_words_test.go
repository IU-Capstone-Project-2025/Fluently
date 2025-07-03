package postgres

import (
	"context"
	"testing"
	"time"

	"fluently/go-backend/internal/repository/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateGetUpdateDeleteLearnedWord(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:        uuid.New(),
		Name:      "Learned User",
		Email:     "learned@example.com",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	word := &models.Word{
		ID:           uuid.New(),
		Word:         "apple",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "яблоко",
		Context:      "I ate an apple",
		AudioURL:     "",
	}
	err = wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	// CREATE
	learned := &models.LearnedWords{
		UserID:           user.ID,
		WordID:           word.ID,
		LearnedAt:        time.Now().Add(-24 * time.Hour),
		LastReviewed:     time.Now(),
		CountOfRevisions: 2,
		ConfidenceScore:  80,
	}
	err = learnedWordRepo.Create(ctx, learned)
	assert.NoError(t, err)

	// GET
	got, err := learnedWordRepo.GetByUserWordID(ctx, user.ID, word.ID)
	assert.NoError(t, err)
	assert.Equal(t, learned.UserID, got.UserID)
	assert.Equal(t, learned.WordID, got.WordID)
	assert.Equal(t, learned.ConfidenceScore, got.ConfidenceScore)

	// UPDATE
	got.ConfidenceScore = 95
	got.CountOfRevisions = 3
	err = learnedWordRepo.Update(ctx, got)
	assert.NoError(t, err)

	updated, err := learnedWordRepo.GetByUserWordID(ctx, user.ID, word.ID)
	assert.NoError(t, err)
	assert.Equal(t, 95, updated.ConfidenceScore)
	assert.Equal(t, 3, updated.CountOfRevisions)

	// LIST
	list, err := learnedWordRepo.ListByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// DELETE
	err = learnedWordRepo.Delete(ctx, user.ID, word.ID)
	assert.NoError(t, err)

	_, err = learnedWordRepo.GetByUserWordID(ctx, user.ID, word.ID)
	assert.Error(t, err) // should return record not found
}
