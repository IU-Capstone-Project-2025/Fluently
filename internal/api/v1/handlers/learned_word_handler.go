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

type LearnedWordHandler struct {
	Repo *postgres.LearnedWordRepository
}

func (h *LearnedWordHandler) GetLearnedWords(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	words, err := h.Repo.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch learned words", http.StatusInternalServerError)
		return
	}

	var resp []schemas.LearenedWordResponse
	for _, w := range words {
		resp = append(resp, schemas.LearenedWordResponse{
			UserID:          w.UserID,
			WordID:          w.WordID,
			LearnedAt:       w.LearnedAt,
			LastReviewed:    w.LastReviewed,
			CntReviewedAt:   w.CountOfRevisions,
			ConfidenceScore: w.ConfidenceScore,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *LearnedWordHandler) GetLearnedWord(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	wordIDStr := chi.URLParam(r, "word_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invelid user_id", http.StatusBadRequest)
		return
	}

	wordID, err := uuid.Parse(wordIDStr)
	if err != nil {
		http.Error(w, "invelid word_id", http.StatusBadRequest)
		return
	}

	word, err := h.Repo.GetByUserWordID(r.Context(), userID, wordID)
	if err != nil {
		http.Error(w, "learned word not found", http.StatusNotFound)
		return
	}

	resp := schemas.LearenedWordResponse{
		UserID:          word.UserID,
		WordID:          word.WordID,
		LearnedAt:       word.LearnedAt,
		LastReviewed:    word.LastReviewed,
		CntReviewedAt:   word.CountOfRevisions,
		ConfidenceScore: word.ConfidenceScore,
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *LearnedWordHandler) CreateLearnedWord(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateLearnedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	word := models.LearnedWords{
		UserID:           req.UserID,
		WordID:           req.WordID,
		LearnedAt:        req.LearnedAt,
		LastReviewed:     req.LastReviewed,
		CountOfRevisions: req.CntReviewed,
		ConfidenceScore:  req.ConfidenceScore,
	}

	if err := h.Repo.Create(r.Context(), &word); err != nil {
		http.Error(w, "failed to create", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LearnedWordHandler) UpdateLearnedWord(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	wordIDStr := chi.URLParam(r, "word_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invelid user_id", http.StatusBadRequest)
		return
	}

	wordID, err := uuid.Parse(wordIDStr)
	if err != nil {
		http.Error(w, "invelid word_id", http.StatusBadRequest)
		return
	}

	var req schemas.CreateLearnedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	word, err := h.Repo.GetByUserWordID(r.Context(), userID, wordID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	word.LearnedAt = req.LearnedAt
	word.LastReviewed = req.LastReviewed
	word.CountOfRevisions = req.CntReviewed
	word.ConfidenceScore = req.ConfidenceScore

	if err := h.Repo.Update(r.Context(), word); err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LearnedWordHandler) DeleteLearnedWord(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	wordIDStr := chi.URLParam(r, "word_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invelid user_id", http.StatusBadRequest)
		return
	}

	wordID, err := uuid.Parse(wordIDStr)
	if err != nil {
		http.Error(w, "invelid word_id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), userID, wordID); err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
