package handler

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SentenceHandler struct {
	Repo *postgres.SentenceRepository
}

func (h *SentenceHandler) ListSentences(w http.ResponseWriter, r *http.Request) {
	wordIDStr := chi.URLParam(r, "word_id")
	wordID, err := uuid.Parse(wordIDStr)
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	sentences, err := h.Repo.ListByWord(r.Context(), wordID)
	if err != nil {
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}

	var resp []schemas.SentenceResponse
	for _, s := range sentences {
		resp = append(resp, schemas.SentenceResponse{
			ID:          s.ID,
			WordID:      s.WordID,
			Sentence:    s.Sentence,
			Translation: s.Translation,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *SentenceHandler) CreateSentence(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateSentenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	s := models.Sentence{
		WordID:      req.WordID,
		Sentence:    req.Sentence,
		Translation: req.Translation,
	}

	if err := h.Repo.Create(r.Context(), &s); err != nil {
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SentenceHandler) UpdateSentence(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	var req schemas.CreateSentenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	sentence, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sentence.WordID = req.WordID
	sentence.Sentence = req.Sentence
	sentence.Translation = req.Translation

	if err := h.Repo.Update(r.Context(), sentence); err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SentenceHandler) DeleteSentence(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
