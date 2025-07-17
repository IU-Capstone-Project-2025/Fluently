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

// TestCreatePickOption tests the creation of a pick option
func TestCreatePickOption(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	wordID := uuid.New()
	sentenceID := uuid.New()

	err := wordRepo.Create(context.Background(), &models.Word{ID: wordID, Word: "example"})
	assert.NoError(t, err)
	err = sentenceRepo.Create(context.Background(), &models.Sentence{ID: sentenceID, WordID: wordID, Sentence: "Example sentence"})
	assert.NoError(t, err)

	body := map[string]interface{}{
		"word_id":     wordID.String(),
		"sentence_id": sentenceID.String(),
		"options":     []string{"one", "two", "three"},
	}

	resp := e.POST("/pick-options").
		WithJSON(body).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.Equal(t, wordID.String(), resp.Value("word_id").String().Raw())
	assert.Equal(t, sentenceID.String(), resp.Value("sentence_id").String().Raw())
	assert.Equal(t, []interface{}{"one", "two", "three"}, resp.Value("options").Array().Raw())
}

// TestGetPickOption tests the retrieval of a pick option
func TestGetPickOption(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	option := models.PickOption{
		ID:         uuid.New(),
		WordID:     uuid.New(),
		SentenceID: uuid.New(),
		Option:     []string{"a", "b", "c"},
	}
	err := pickOptionRepo.Create(context.Background(), &option)
	assert.NoError(t, err)

	resp := e.GET("/pick-options/" + option.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, option.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, []interface{}{"a", "b", "c"}, resp.Value("options").Array().Raw())
}

// TestUpdatePickOption tests the update of a pick option
func TestUpdatePickOption(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	option := models.PickOption{
		ID:         uuid.New(),
		WordID:     uuid.New(),
		SentenceID: uuid.New(),
		Option:     []string{"x", "y", "z"},
	}
	err := pickOptionRepo.Create(context.Background(), &option)
	assert.NoError(t, err)

	update := map[string]interface{}{
		"word_id":     option.WordID.String(),
		"sentence_id": option.SentenceID.String(),
		"options":     []string{"new1", "new2", "new3"},
	}

	e.PUT("/pick-options/" + option.ID.String()).
		WithJSON(update).
		Expect().
		Status(http.StatusOK)

	updated, err := pickOptionRepo.GetByID(context.Background(), option.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"new1", "new2", "new3"}, updated.Option)
}

// TestDeletePickOption tests the deletion of a pick option
func TestDeletePickOption(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	option := models.PickOption{
		ID:         uuid.New(),
		WordID:     uuid.New(),
		SentenceID: uuid.New(),
		Option:     []string{"del1", "del2", "del3"},
	}
	err := pickOptionRepo.Create(context.Background(), &option)
	assert.NoError(t, err)

	e.DELETE("/pick-options/" + option.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	_, err = pickOptionRepo.GetByID(context.Background(), option.ID)
	assert.Error(t, err)
}

// TestListPickOptions tests the listing of pick options
func TestListPickOptions(t *testing.T) {
	setupTest(t)
	e := httpexpect.Default(t, testServer.URL)

	wordID := uuid.New()
	sentenceID := uuid.New()

	err := wordRepo.Create(context.Background(), &models.Word{ID: wordID, Word: "ListWord"})
	assert.NoError(t, err)

	options := []models.PickOption{
		{
			ID:         uuid.New(),
			WordID:     wordID,
			SentenceID: sentenceID,
			Option:     []string{"1", "2", "3"},
		},
		{
			ID:         uuid.New(),
			WordID:     wordID,
			SentenceID: sentenceID,
			Option:     []string{"a", "b", "c"},
		},
	}

	for _, o := range options {
		err := pickOptionRepo.Create(context.Background(), &o)
		assert.NoError(t, err)
	}

	resp := e.GET("/words/" + wordID.String() + "/pick-options").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	assert.Len(t, resp.Raw(), 2)
}
