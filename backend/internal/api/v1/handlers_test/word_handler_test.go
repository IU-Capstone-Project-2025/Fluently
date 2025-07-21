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

// TestCreateWord tests the creation of a new word
func TestCreateWord(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	req := map[string]interface{}{
		"word":           "apple",
		"cefr_level":     "A1",
		"part_of_speech": "noun",
		"translation":    "яблоко",
		"context":        "I ate an apple",
		"audio_url":      "http://tyanka-vovanka.com",
	}

	resp := e.POST("/api/v1/words").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.Equal(t, "apple", resp.Value("word").String().Raw())
	assert.Equal(t, "A1", resp.Value("cefr_level").String().Raw())
	assert.Equal(t, "noun", resp.Value("part_of_speech").String().Raw())
	assert.Equal(t, "яблоко", resp.Value("translation").String().Raw())
	assert.Equal(t, "I ate an apple", resp.Value("context").String().Raw())
	assert.Equal(t, "http://tyanka-vovanka.com", resp.Value("audio_url").String().Raw())
	assert.NotEmpty(t, resp.Value("id").String().Raw())
}

// TestListWords tests the retrieval of all words
func TestListWords(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:           uuid.New(),
		Word:         "banana",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "банан",
		Context:      "I like banana",
		AudioURL:     "http://banana.audio",
	}

	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/words").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	length := int(resp.Length().Raw())
	found := false
	for i := 0; i < length; i++ {
		obj := resp.Value(i).Object()
		if obj.Value("word").String().Raw() == "banana" {
			found = true
			assert.Equal(t, "A1", obj.Value("cefr_level").String().Raw())
			assert.Equal(t, "noun", obj.Value("part_of_speech").String().Raw())
			assert.Equal(t, "банан", obj.Value("translation").String().Raw())
		}
	}
	assert.True(t, found)
}

// TestGetWord tests the retrieval of a word
func TestGetWord(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:           uuid.New(),
		Word:         "cherry",
		CEFRLevel:    "B1",
		PartOfSpeech: "noun",
		Translation:  "вишня",
		Context:      "I ate a cherry",
		AudioURL:     "http://cherry.audio",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/words/" + word.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, "cherry", resp.Value("word").String().Raw())
	assert.Equal(t, "B1", resp.Value("cefr_level").String().Raw())
	assert.Equal(t, "noun", resp.Value("part_of_speech").String().Raw())
	assert.Equal(t, "вишня", resp.Value("translation").String().Raw())
}

// TestUpdateWord tests the update of a word
func TestUpdateWord(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:           uuid.New(),
		Word:         "orange",
		CEFRLevel:    "A2",
		PartOfSpeech: "noun",
		Translation:  "апельсин",
		AudioURL:     "http://orange.audio",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"word":           "orange-updated",
		"cefr_level":     "B2",
		"part_of_speech": "noun",
		"translation":    "апельсин обновленный",
		"context":        "Updated context",
		"audio_url":      "http://updated.audio",
	}

	resp := e.PUT("/api/v1/words/" + word.ID.String()).
		WithJSON(updateBody).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, "orange-updated", resp.Value("word").String().Raw())
	assert.Equal(t, "B2", resp.Value("cefr_level").String().Raw())
	assert.Equal(t, "апельсин обновленный", resp.Value("translation").String().Raw())
	assert.Equal(t, "Updated context", resp.Value("context").String().Raw())
	assert.Equal(t, "http://updated.audio", resp.Value("audio_url").String().Raw())
}

// TestDeleteWord tests the deletion of a word
func TestDeleteWord(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	word := models.Word{
		ID:           uuid.New(),
		Word:         "pear",
		CEFRLevel:    "A2",
		PartOfSpeech: "noun",
		Translation:  "груша",
	}
	err := wordRepo.Create(context.Background(), &word)
	assert.NoError(t, err)

	e.DELETE("/api/v1/words/" + word.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	e.GET("/api/v1/words/" + word.ID.String()).
		Expect().
		Status(http.StatusNotFound)
}
