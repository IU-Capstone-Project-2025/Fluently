package handlers

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
)

type TopicHandler struct {
	Repo *postgres.TopicRepository
}

func buildTopicResponse(topic *models.Topic) schemas.TopicResponse {
	return schemas.TopicResponse{
		ID:    topic.ID.String(),
		Title: topic.Title,
	}
}

// GetTopic godoc
// @Summary      Get topic by ID
// @Description  Returns a topic by its unique identifier
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Topic ID"
// @Success      200  {object}  schemas.TopicResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Router       /api/v1/topics/{id} [get]
func (h *TopicHandler) GetTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	topic, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "topic not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildTopicResponse(topic))
}

// GetMainTopic godoc
// @Summary      Get main topic by ID
// @Description  Returns the main topic (root parent) for a given topic ID
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Topic ID"
// @Success      200  {object}  schemas.TopicResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/topics/{id}/main [get]
func (h *TopicHandler) GetMainTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	topic, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "topic not found", http.StatusNotFound)
		return
	}

	for topic.ParentID != nil {
		topic, err = h.Repo.GetByID(r.Context(), *topic.ParentID)
		if err != nil {
			http.Error(w, "failed to fetch parent topic", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildTopicResponse(topic))
}

// GetPathToMainTopic godoc
// @Summary      Get path to main topic
// @Description  Returns the path from a topic to its main topic (root parent)
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Topic ID"
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/topics/{id}/path [get]
func (h *TopicHandler) GetPathToMainTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var path []uuid.UUID

	topic, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "topic not found", http.StatusNotFound)
		return
	}
	path = append(path, topic.ID)

	for topic.ParentID != nil {
		topic, err = h.Repo.GetByID(r.Context(), *topic.ParentID)
		if err != nil {
			http.Error(w, "failed to fetch parent topic", http.StatusInternalServerError)
			return
		}
		path = append(path, topic.ID)
	}

	strPath := make([]string, len(path))
	for i, id := range path {
		strPath[i] = id.String()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"topics": strPath,
	})
}

// CreateTopic godoc
// @Summary      Create a new topic
// @Description  Adds a new topic
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        topic  body      schemas.CreateTopicRequest  true  "Topic data"
// @Success      201  {object}  schemas.TopicResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/topics/ [post]
func (h *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateTopicRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	topic := models.Topic{
		Title: req.Title,
	}

	if err := h.Repo.Create(r.Context(), &topic); err != nil {
		http.Error(w, "failed to create topic", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildTopicResponse(&topic))
}

// UpdateTopic godoc
// @Summary      Update a topic
// @Description  Updates an existing topic by ID
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                     true  "Topic ID"
// @Param        topic  body      schemas.CreateTopicRequest true  "Topic data"
// @Success      200  {object}  schemas.TopicResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/topics/{id} [put]
func (h *TopicHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req schemas.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	topic, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "topic not found", http.StatusNotFound)
		return
	}

	topic.Title = req.Title

	if err := h.Repo.Update(r.Context(), topic); err != nil {
		http.Error(w, "failed to update topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildTopicResponse(topic))
}

// DeleteTopic godoc
// @Summary      Delete a topic
// @Description  Deletes a topic by ID
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Topic ID"
// @Success      204  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/topics/{id} [delete]
func (h *TopicHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
