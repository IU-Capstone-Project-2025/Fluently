package postgres

import (
	"context"
	"testing"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetSentence(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		ID:           uuid.New(),
		Word:         "run",
		Translation:  "беги",
		PartOfSpeech: "verb",
		CEFRLevel:    "A1",
		Context:      "I run every day.",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	sentence := &models.Sentence{
		WordID:      word.ID,
		Sentence:    "I like to run in the park.",
		Translation: "Мне нравится бегать в парке",
		AudioURL:    "http://jojoref.mp3",
	}

	err = sentenceRepo.Create(ctx, sentence)
	assert.NoError(t, err)

	found, err := sentenceRepo.GetByID(ctx, sentence.ID)
	assert.NoError(t, err)
	assert.Equal(t, sentence.Sentence, found.Sentence)
	assert.Equal(t, sentence.Translation, found.Translation)
	assert.Equal(t, sentence.AudioURL, found.AudioURL)
}

func TestUpdateSentence(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		ID:           uuid.New(),
		Word:         "jump",
		Translation:  "прыгать",
		PartOfSpeech: "verb",
		CEFRLevel:    "A1",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	sentence := &models.Sentence{
		WordID:      word.ID,
		Sentence:    "I can jump high.",
		Translation: "Я могу прыгать высоко.",
	}

	err = sentenceRepo.Create(ctx, sentence)
	assert.NoError(t, err)

	sentence.Translation = "Я умею прыгать высоко."
	sentence.AudioURL = "http://nojojoref.mp3"
	err = sentenceRepo.Update(ctx, sentence)
	assert.NoError(t, err)

	updated, err := sentenceRepo.GetByID(ctx, sentence.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Я умею прыгать высоко.", updated.Translation)
	assert.Equal(t, "http://nojojoref.mp3", updated.AudioURL)
}

func TestDeleteSentence(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		ID:           uuid.New(),
		Word:         "sleep",
		Translation:  "спать",
		PartOfSpeech: "verb",
		CEFRLevel:    "A1",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	sentence := &models.Sentence{
		WordID:      word.ID,
		Sentence:    "I sleep every second.",
		Translation: "Я сплю каждую секунду.",
	}

	err = sentenceRepo.Create(ctx, sentence)
	assert.NoError(t, err)

	err = sentenceRepo.Delete(ctx, sentence.ID)
	assert.NoError(t, err)

	_, err = sentenceRepo.GetByID(ctx, sentence.ID)
	assert.Error(t, err) // "record not found"
}
