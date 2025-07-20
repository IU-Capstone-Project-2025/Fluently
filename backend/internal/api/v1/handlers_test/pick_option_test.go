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

// TestCreatePickOption tests the creation of a new pick option
func TestCreatePickOption(t *testing.T) {
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
		Sentence:    "Test sentence",
		Translation: "Тестовое предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	req := map[string]interface{}{
		"word_id":     word.ID.String(),
		"sentence_id": sentence.ID.String(),
		"options":     []string{"one", "two", "three"},
	}

	resp := e.POST("/api/v1/pick-options").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.Equal(t, word.ID.String(), resp.Value("word_id").String().Raw())
	assert.Equal(t, sentence.ID.String(), resp.Value("sentence_id").String().Raw())
	assert.Equal(t, []interface{}{"one", "two", "three"}, resp.Value("options").Array().Raw())
}

// TestGetPickOption tests the retrieval of a pick option
func TestGetPickOption(t *testing.T) {
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
		Sentence:    "Test sentence",
		Translation: "Тестовое предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	pickOption := models.PickOption{
		ID:         uuid.New(),
		WordID:     word.ID,
		SentenceID: sentence.ID,
		Option:     models.StringArray{"a", "b", "c"},
	}
	err = pickOptionRepo.Create(context.Background(), &pickOption)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/pick-options/" + pickOption.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, pickOption.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, []interface{}{"a", "b", "c"}, resp.Value("options").Array().Raw())
}

// TestUpdatePickOption tests the update of a pick option
func TestUpdatePickOption(t *testing.T) {
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
		Sentence:    "Test sentence",
		Translation: "Тестовое предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	pickOption := models.PickOption{
		ID:         uuid.New(),
		WordID:     word.ID,
		SentenceID: sentence.ID,
		Option:     models.StringArray{"x", "y", "z"},
	}
	err = pickOptionRepo.Create(context.Background(), &pickOption)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"word_id":     word.ID.String(),
		"sentence_id": sentence.ID.String(),
		"options":     []string{"new1", "new2", "new3"},
	}

	resp := e.PUT("/api/v1/pick-options/" + pickOption.ID.String()).
		WithJSON(updateBody).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Check that the options were updated
	updatedOptions := resp.Value("options").Array()
	assert.Equal(t, 3, int(updatedOptions.Length().Raw()))
	assert.Equal(t, "new1", updatedOptions.Value(0).String().Raw())
	assert.Equal(t, "new2", updatedOptions.Value(1).String().Raw())
	assert.Equal(t, "new3", updatedOptions.Value(2).String().Raw())
}

// TestDeletePickOption tests the deletion of a pick option
func TestDeletePickOption(t *testing.T) {
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
		Sentence:    "Test sentence",
		Translation: "Тестовое предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	pickOption := models.PickOption{
		ID:         uuid.New(),
		WordID:     word.ID,
		SentenceID: sentence.ID,
		Option:     models.StringArray{"delete", "me", "now"},
	}
	err = pickOptionRepo.Create(context.Background(), &pickOption)
	assert.NoError(t, err)

	e.DELETE("/api/v1/pick-options/" + pickOption.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	// Verify it was deleted
	_, err = pickOptionRepo.GetByID(context.Background(), pickOption.ID)
	assert.Error(t, err)
}

// TestListPickOptions tests the listing of pick options for a word
func TestListPickOptions(t *testing.T) {
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
		Sentence:    "Test sentence",
		Translation: "Тестовое предложение",
	}
	err = sentenceRepo.Create(context.Background(), &sentence)
	assert.NoError(t, err)

	// Create two pick options for the same word
	pickOption1 := models.PickOption{
		ID:         uuid.New(),
		WordID:     word.ID,
		SentenceID: sentence.ID,
		Option:     models.StringArray{"a", "b", "c"},
	}
	err = pickOptionRepo.Create(context.Background(), &pickOption1)
	assert.NoError(t, err)

	pickOption2 := models.PickOption{
		ID:         uuid.New(),
		WordID:     word.ID,
		SentenceID: sentence.ID,
		Option:     models.StringArray{"d", "e", "f"},
	}
	err = pickOptionRepo.Create(context.Background(), &pickOption2)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/words/" + word.ID.String() + "/pick-options").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	// Should have 2 pick options
	assert.Equal(t, 2, int(resp.Length().Raw()))

	// Check that both pick options are present
	found1 := false
	found2 := false
	for i := 0; i < int(resp.Length().Raw()); i++ {
		option := resp.Value(i).Object()
		id := option.Value("id").String().Raw()
		if id == pickOption1.ID.String() {
			found1 = true
		}
		if id == pickOption2.ID.String() {
			found2 = true
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)
}
