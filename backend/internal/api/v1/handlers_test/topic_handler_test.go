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

	e := httpexpect.Default(t, testServer.URL)

	req := map[string]interface{}{
		"title":     "New Topic",
		"parent_id": nil,
	}

	resp := e.POST("/api/v1/topics").
		WithJSON(req).
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

	resp := e.GET("/api/v1/topics/" + topic.ID.String()).
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
		Title: "Original Title",
	}
	err := topicRepo.Create(context.Background(), &topic)
	assert.NoError(t, err)

	updateBody := map[string]interface{}{
		"title": "Updated Title",
	}

	resp := e.PUT("/api/v1/topics/" + topic.ID.String()).
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
		Title: "Delete Me",
	}
	err := topicRepo.Create(context.Background(), &topic)
	assert.NoError(t, err)

	e.DELETE("/api/v1/topics/" + topic.ID.String()).
		Expect().
		Status(http.StatusNoContent)

	e.GET("/api/v1/topics/" + topic.ID.String()).
		Expect().
		Status(http.StatusNotFound)
}

// TestGetMainTopic tests the retrieval of a main topic
func TestGetMainTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	// Create a root topic
	rootTopic := models.Topic{
		ID:    uuid.New(),
		Title: "Root Topic",
	}
	err := topicRepo.Create(context.Background(), &rootTopic)
	assert.NoError(t, err)

	// Create a child topic
	childTopic := models.Topic{
		ID:       uuid.New(),
		Title:    "Child Topic",
		ParentID: &rootTopic.ID,
	}
	err = topicRepo.Create(context.Background(), &childTopic)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/topics/root-topic/" + childTopic.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	assert.Equal(t, rootTopic.ID.String(), resp.Value("id").String().Raw())
	assert.Equal(t, "Root Topic", resp.Value("title").String().Raw())
}

// TestGetPathToMainTopic tests the retrieval of path to main topic
func TestGetPathToMainTopic(t *testing.T) {
	setupTest(t)

	e := httpexpect.Default(t, testServer.URL)

	// Create a root topic
	rootTopic := models.Topic{
		ID:    uuid.New(),
		Title: "Root Topic",
	}
	err := topicRepo.Create(context.Background(), &rootTopic)
	assert.NoError(t, err)

	// Create a child topic
	childTopic := models.Topic{
		ID:       uuid.New(),
		Title:    "Child Topic",
		ParentID: &rootTopic.ID,
	}
	err = topicRepo.Create(context.Background(), &childTopic)
	assert.NoError(t, err)

	resp := e.GET("/api/v1/topics/path-to-root/" + childTopic.ID.String()).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// The response should have a "topics" field containing an array
	topicsArray := resp.Value("topics").Array()

	// Should have 2 topics in the path: child and root
	assert.Equal(t, 2, int(topicsArray.Length().Raw()))

	// First should be the child topic ID
	assert.Equal(t, childTopic.ID.String(), topicsArray.Value(0).String().Raw())

	// Second should be the root topic ID
	assert.Equal(t, rootTopic.ID.String(), topicsArray.Value(1).String().Raw())
}
