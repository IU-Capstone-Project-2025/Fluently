package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
)

// LearnedWordHandler is a handler for learned words
type LearnedWordHandler struct {
	Repo *postgres.LearnedWordRepository
}

// buildLearnedWordResponse builds a response from a learned word
func buildLearnedWordResponse(word *models.LearnedWords) schemas.LearnedWordResponse {
	return schemas.LearnedWordResponse{
		UserID:          word.UserID,
		WordID:          word.WordID,
		LearnedAt:       word.LearnedAt,
		LastReviewed:    word.LastReviewed,
		CntReviewedAt:   word.CountOfRevisions,
		ConfidenceScore: word.ConfidenceScore,
	}
}

// GetLearnedWords returns all learned words
func (h *LearnedWordHandler) GetLearnedWords(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/learned-words"
	method := r.Method
	statusCode := 200
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		statusCode = 401
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	words, err := h.Repo.ListByUserID(r.Context(), user.ID)
	if err != nil {
		statusCode = 500
		http.Error(w, "failed to fetch learned words", http.StatusInternalServerError)
		return
	}

	var resp []schemas.LearnedWordResponse
	for _, w := range words {
		resp = append(resp, buildLearnedWordResponse(&w))
	}

	// Return the list of learned words
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetLearnedWord returns a single learned word for a user
func (h *LearnedWordHandler) GetLearnedWord(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/learned-words/{word_id}"
	method := r.Method
	statusCode := 200
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		statusCode = 401
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	wordID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		statusCode = 400
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	word, err := h.Repo.GetByUserWordID(r.Context(), user.ID, wordID)
	if err != nil {
		statusCode = 404
		http.Error(w, "learned word not found", http.StatusNotFound)
		return
	}

	// Return the learned word
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildLearnedWordResponse(word))
}

// CreateLearnedWord creates a new learned word
func (h *LearnedWordHandler) CreateLearnedWord(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/learned-words"
	method := r.Method
	statusCode := 201
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		statusCode = 401
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req schemas.CreateLearnedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		statusCode = 400
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID != uuid.Nil && req.UserID != user.ID {
		statusCode = 400
		http.Error(w, "user_id mismatch", http.StatusBadRequest)
		return
	}

	word := models.LearnedWords{
		UserID:           user.ID,
		WordID:           req.WordID,
		LearnedAt:        req.LearnedAt,
		LastReviewed:     req.LastReviewed,
		CountOfRevisions: req.CntReviewed,
		ConfidenceScore:  req.ConfidenceScore,
	}

	if err := h.Repo.Create(r.Context(), &word); err != nil {
		statusCode = 500
		http.Error(w, "failed to create learned word", http.StatusInternalServerError)
		return
	}

	// Return the created learned word
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// UpdateLearnedWord updates a learned word
func (h *LearnedWordHandler) UpdateLearnedWord(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/learned-words/{word_id}"
	method := r.Method
	statusCode := 200
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		statusCode = 401
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req schemas.CreateLearnedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		statusCode = 400
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	word, err := h.Repo.GetByUserWordID(r.Context(), user.ID, req.WordID)
	if err != nil {
		statusCode = 404
		http.Error(w, "learned word not found", http.StatusNotFound)
		return
	}

	word.LearnedAt = req.LearnedAt
	word.LastReviewed = req.LastReviewed
	word.CountOfRevisions = req.CntReviewed
	word.ConfidenceScore = req.ConfidenceScore

	if err := h.Repo.Update(r.Context(), word); err != nil {
		statusCode = 500
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	// Return the updated learned word
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteLearnedWord deletes a learned word
func (h *LearnedWordHandler) DeleteLearnedWord(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/learned-words/{word_id}"
	method := r.Method
	statusCode := 204
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()
	wordID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		statusCode = 400
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		statusCode = 401
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.Repo.Delete(r.Context(), user.ID, wordID); err != nil {
		statusCode = 500
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	// Return no content
	w.WriteHeader(http.StatusNoContent)
}
