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

func (h *TopicHandler) GetTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.Repo.ListTopics(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch topics", http.StatusBadRequest)
		return
	}

	var resp []schemas.TopicResponse
	for _, t := range topics {
		resp = append(resp, buildTopicResponse(&t))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *TopicHandler) GetTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
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

func (h *TopicHandler) GetMainTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
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

func (h *TopicHandler) GetPathToMainTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(strPath)
}

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

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildTopicResponse(&topic))
}

func (h *TopicHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
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
}

func (h *TopicHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
