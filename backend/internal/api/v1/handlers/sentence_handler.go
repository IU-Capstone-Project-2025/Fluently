package handlers

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
)

type SentenceHandler struct {
	Repo *postgres.SentenceRepository
}

func buildSentenceResponse(sentence *models.Sentence) schemas.SentenceResponse {
	return schemas.SentenceResponse{
		ID:          sentence.ID.String(),
		WordID:      sentence.WordID.String(),
		Sentence:    sentence.Sentence,
		Translation: sentence.Translation,
		AudioURL:    sentence.AudioURL,
	}
}

func (h *SentenceHandler) ListSentences(w http.ResponseWriter, r *http.Request) {
	wordID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	sentences, err := h.Repo.ListByWord(r.Context(), wordID)
	if err != nil {
		http.Error(w, "failed to fetch sentences", http.StatusInternalServerError)
		return
	}

	var resp []schemas.SentenceResponse
	for _, s := range sentences {
		resp = append(resp, buildSentenceResponse(&s))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *SentenceHandler) CreateSentence(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateSentenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	s := models.Sentence{
		WordID:      req.WordID,
		Sentence:    req.Sentence,
		Translation: req.Translation,
		AudioURL:    req.AudioURL,
	}

	if err := h.Repo.Create(r.Context(), &s); err != nil {
		http.Error(w, "failed to create sentence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildSentenceResponse(&s))
}

func (h *SentenceHandler) UpdateSentence(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req schemas.CreateSentenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sentence, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "sentence not found", http.StatusNotFound)
		return
	}

	sentence.WordID = req.WordID
	sentence.Sentence = req.Sentence
	sentence.Translation = req.Translation

	if req.AudioURL != "" {
		sentence.AudioURL = req.AudioURL
	}

	if err := h.Repo.Update(r.Context(), sentence); err != nil {
		http.Error(w, "failed to update sentence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildSentenceResponse(sentence))
}

func (h *SentenceHandler) DeleteSentence(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete sentence", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
