package postgres

import (
	"context"
	"testing"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetWord(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		ID:           uuid.New(),
		Word:         "apple",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "яблоко",
		Context:      "I ate an apple",
		AudioURL:     "",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	found, err := wordRepo.GetByID(ctx, word.ID)
	assert.NoError(t, err)
	assert.Equal(t, word.Word, found.Word)
	assert.Equal(t, word.Translation, found.Translation)
	assert.Equal(t, word.PartOfSpeech, found.PartOfSpeech)
	assert.Equal(t, word.Context, found.Context)
	assert.Equal(t, word.AudioURL, found.AudioURL)
}

func TestListWord(t *testing.T) {
	ctx := context.Background()

	words, err := wordRepo.ListWords(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, words)
}

func TestUpdateWord(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		ID:           uuid.New(),
		Word:         "banana",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "банан",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	word.Translation = "Новый перевод"
	word.AudioURL = "Rickroll URL"

	err = wordRepo.Update(ctx, word)
	assert.NoError(t, err)

	updated, err := wordRepo.GetByID(ctx, word.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Новый перевод", updated.Translation)
	assert.Equal(t, "Rickroll URL", updated.AudioURL)
}

func DeleteWord(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		ID:           uuid.New(),
		Word:         "orange",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	err = wordRepo.Delete(ctx, word.ID)
	assert.NoError(t, err)

	_, err = wordRepo.GetByID(ctx, word.ID)
	assert.Error(t, err) // Not Found
}
