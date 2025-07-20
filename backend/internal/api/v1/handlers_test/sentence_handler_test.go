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

// TestCreateSentence tests the creation of a new sentence
func TestCreateSentence(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:   uuid.New(),
		Word: "test",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	req := map[string]interface{}{
		"word_id":     word.ID.String(),
		"sentence":    "Example sentence",
		"translation": "Пример предложения",
		"audio_url":   "http://example.com/audio.mp3",
	}

	resp := e.POST("/api/v1/sentences").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.NotEmpty(t, resp.Value("id").String().Raw())
	assert.Equal(t, word.ID.String(), resp.Value("word_id").String().Raw())
	assert.Equal(t, "Example sentence", resp.Value("sentence").String().Raw())
	assert.Equal(t, "Пример предложения", resp.Value("translation").String().Raw())
	assert.Equal(t, "http://example.com/audio.mp3", resp.Value("audio_url").String().Raw())
}

// TestListSentences tests the listing of sentences for a word
func TestListSentences(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:   uuid.New(),
		Word: "test",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	// Create two sentences for the same word
	sentence1 := models.Sentence{
		ID:          uuid.New(),
		WordID:      word.ID,
		Sentence:    "First sentence",
		Translation: "Первое предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence1)
	assert.NoError(t, err)

	sentence2 := models.Sentence{
		ID:          uuid.New(),
		WordID:      word.ID,
		Sentence:    "Second sentence",
		Translation: "Второе предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence2)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/words/" + word.ID.String() + "/sentences").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	// Should have 2 sentences
	assert.Equal(t, 2, int(resp.Length().Raw()))

	// Check that both sentences are present
	found1 := false
	found2 := false
	for i := 0; i < int(resp.Length().Raw()); i++ {
		sentence := resp.Value(i).Object()
		id := sentence.Value("id").String().Raw()
		if id == sentence1.ID.String() {
			found1 = true
		}
		if id == sentence2.ID.String() {
			found2 = true
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)
}

// TestUpdateSentence tests the update of a sentence
func TestUpdateSentence(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:   uuid.New(),
		Word: "test",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	sentence := models.Sentence{
		ID:          uuid.New(),
		WordID:      word.ID,
		Sentence:    "Original sentence",
		Translation: "Оригинальное предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"word_id":     word.ID.String(),
		"sentence":    "Updated sentence",
		"translation": "Обновленный перевод",
		"audio_url":   "http://example.com/new_audio.mp3",
	}

	resp := e.PUT("/api/v1/sentences/" + sentence.ID.String()).
		WithJSON(updateBody).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, "Updated sentence", resp.Value("sentence").String().Raw())
	assert.Equal(t, "Обновленный перевод", resp.Value("translation").String().Raw())
	assert.Equal(t, "http://example.com/new_audio.mp3", resp.Value("audio_url").String().Raw())
}

// TestDeleteSentence tests the deletion of a sentence
func TestDeleteSentence(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:   uuid.New(),
		Word: "test",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	sentence := models.Sentence{
		ID:          uuid.New(),
		WordID:      word.ID,
		Sentence:    "Delete me",
		Translation: "Удалить меня",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	e.DELETE("/api/v1/sentences/" + sentence.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	// Verify it was deleted
	_, err = sentenceRepo.GetByID(context.Background(), sentence.ID)
	assert.Error(t, err)
}
