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

// TestCreateSentence tests the creation of a sentence
func TestCreateSentence(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:   uuid.New(),
		Word: "testword",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	req := map[string]interface{}{
		"word_id":     word.ID.String(),
		"sentence":    "Example sentence",
		"translation": "Пример предложения",
		"audio_url":   "http://example.com/audio.mp3",
	}

	resp := e.POST("/sentences").
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

// TestListSentences tests the listing of sentences
func TestListSentences(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:   uuid.New(),
		Word: "testword2",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	s1 := models.Sentence{
		WordID:      word.ID,
		Sentence:    "Sentence 1",
		Translation: "Предложение 1",
	}
	err = sentenceRepo.Create(context.Background(), &s1)
	assert.NoError(t, err)

	s2 := models.Sentence{
		WordID:      word.ID,
		Sentence:    "Sentence 2",
		Translation: "Предложение 2",
		AudioURL:    "http://example.com/audio2.mp3",
	}
	err = sentenceRepo.Create(context.Background(), &s2)
	assert.NoError(t, err)

	resp := e.GET("/words/" + word.ID.String() + "/sentences").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	assert.Equal(t, float64(2), resp.Length().Raw())

	length := int(resp.Length().Raw())
	found := false
	for i := 0; i < length; i++ {
		sentence := resp.Value(i).Object()
		if sentence.Value("sentence").String().Raw() == "Sentence 1" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

// TestUpdateSentence tests the update of a sentence
func TestUpdateSentence(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:   uuid.New(),
		Word: "testword3",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	sentence := models.Sentence{
		WordID:      word.ID,
		Sentence:    "Old sentence",
		Translation: "Старый перевод",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"word_id":     word.ID.String(),
		"sentence":    "Updated sentence",
		"translation": "Обновленный перевод",
		"audio_url":   "http://example.com/new_audio.mp3",
	}

	resp := e.PUT("/sentences/" + sentence.ID.String()).
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
		Word: "testword4",
	}

	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	sentence := models.Sentence{
		WordID:   word.ID,
		Sentence: "To delete",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	e.DELETE("/sentences/" + sentence.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	_, err = sentenceRepo.GetByID(context.Background(), sentence.ID)
	assert.Error(t, err)
}
