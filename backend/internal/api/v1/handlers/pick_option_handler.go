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

type PickOptionHandler struct {
	Repo *postgres.PickOptionRepository
}

func buildPickOptionResponse(option *models.PickOption) schemas.PickOptionResponse {
	return schemas.PickOptionResponse{
		ID:         option.ID.String(),
		WordID:     option.WordID.String(),
		SentenceID: option.SentenceID.String(),
		Options:    option.Option,
	}
}

func (h *PickOptionHandler) CreatePickOption(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreatePickOptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	wordID, err := uuid.Parse(req.WordID)
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	sentenceID, err := uuid.Parse(req.SentenceID)
	if err != nil {
		http.Error(w, "invalid sentence_id", http.StatusBadRequest)
		return
	}

	option := models.PickOption{
		ID:         uuid.New(),
		WordID:     wordID,
		SentenceID: sentenceID,
		Option:     req.Options,
	}

	if err := h.Repo.Create(r.Context(), &option); err != nil {
		http.Error(w, "failed to create option", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPickOptionResponse(&option))
}

func (h *PickOptionHandler) UpdatePickOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	var req schemas.CreatePickOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	option, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	}

	option.Option = req.Options

	if err := h.Repo.Update(r.Context(), option); err != nil {
		http.Error(w, "failed to update topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PickOptionHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete option", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PickOptionHandler) GetTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	option, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPickOptionResponse(option))
}

func (h *PickOptionHandler) List(w http.ResponseWriter, r *http.Request) {
	options, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch options", http.StatusBadRequest)
		return
	}

	var resp []schemas.PickOptionResponse
	for _, o := range options {
		resp = append(resp, buildPickOptionResponse(&o))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
