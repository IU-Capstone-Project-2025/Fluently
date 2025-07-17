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

// TestCreateTopic tests the creation of a new topic
func TestCreateTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  testServer.URL,
		Client:   testServer.Client(),
		Reporter: httpexpect.NewAssertReporter(t),
	})

	reqBody := map[string]interface{}{
		"title": "New Topic",
	}

	resp := e.POST("/topics/").
		WithJSON(reqBody).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	assert.Equal(t, "New Topic", resp.Value("title").String().Raw())
	assert.NotEmpty(t, resp.Value("id").String().Raw())
}

// TestGetTopic tests the retrieval of a topic
func TestGetTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	topic := models.Topic{
		ID:    uuid.New(),
		Title: "Test Topic",
	}
	err := topicRepo.Create(context.Background(), &topic)
	assert.NoError(t, err)

	resp := e.GET("/topics/" + topic.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, topic.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, "Test Topic", resp.Value("title").String().Raw())
}

// TestUpdateTopic tests the update of a topic
func TestUpdateTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	topic := models.Topic{
		ID:    uuid.New(),
		Title: "Old Title",
	}
	err := topicRepo.Create(context.Background(), &topic)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"title": "Updated Title",
	}

	resp := e.PUT("/topics/" + topic.ID.String()).
		WithJSON(updateBody).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, topic.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, "Updated Title", resp.Value("title").String().Raw())
}

// TestDeleteTopic tests the deletion of a topic
func TestDeleteTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	topic := models.Topic{
		ID:    uuid.New(),
		Title: "Topic to delete",
	}
	err := topicRepo.Create(context.Background(), &topic)
	assert.NoError(t, err)

	e.DELETE("/topics/" + topic.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	e.GET("/topics/" + topic.ID.String()).
		Expect().
		Status(http.StatusNotFound)
}

// TestGetMainTopic tests the retrieval of the main topic
func TestGetMainTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	rootTopic := models.Topic{
		ID:    uuid.New(),
		Title: "Root Topic",
	}
	err := topicRepo.Create(context.Background(), &rootTopic)
	assert.NoError(t, err)

	childTopic := models.Topic{
		ID:       uuid.New(),
		Title:    "Child Topic",
		ParentID: &rootTopic.ID,
	}
	err = topicRepo.Create(context.Background(), &childTopic)
	assert.NoError(t, err)

	resp := e.GET("/topics/root-topic/" + childTopic.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, rootTopic.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, "Root Topic", resp.Value("title").String().Raw())
}

// TestGetPathToMainTopic tests the retrieval of the path to the main topic
func TestGetPathToMainTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	rootTopic := models.Topic{
		ID:    uuid.New(),
		Title: "Root Topic",
	}
	err := topicRepo.Create(context.Background(), &rootTopic)
	assert.NoError(t, err)

	childTopic := models.Topic{
		ID:       uuid.New(),
		Title:    "Child Topic",
		ParentID: &rootTopic.ID,
	}
	err = topicRepo.Create(context.Background(), &childTopic)
	assert.NoError(t, err)

	resp := e.GET("/topics/path-to-root/" + childTopic.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	topics := resp.Value("topics").Array()

	assert.Equal(t, childTopic.ID.String(), topics.Value(0).String().Raw())
	assert.Equal(t, rootTopic.ID.String(), topics.Value(1).String().Raw())
}
