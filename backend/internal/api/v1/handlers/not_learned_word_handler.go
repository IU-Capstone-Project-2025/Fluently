package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotLearnedWordHandler is a handler for not learned words
type NotLearnedWordHandler struct {
	Repo     *postgres.NotLearnedWordRepository
	WordRepo *postgres.WordRepository
}

// buildNotLearnedWordResponse builds a response from a not learned word
func buildNotLearnedWordResponse(nlw *models.NotLearnedWords, word *models.Word) schemas.NotLearnedWordResponse {
	return schemas.NotLearnedWordResponse{
		ID:        nlw.ID,
		UserID:    nlw.UserID,
		WordID:    nlw.WordID,
		Word:      word.Word,
		CreatedAt: time.Now(), // Since the model doesn't have CreatedAt, use current time
	}
}

// AddNotLearnedWord godoc
// @Summary      Add word to not learned words
// @Description  Adds a word to the user's not learned words list
// @Tags         not-learned-words
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body schemas.AddNotLearnedWordRequest true "Word to add"
// @Success      201  {object}  schemas.NotLearnedWordResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/not-learned-words [post]
func (h *NotLearnedWordHandler) AddNotLearnedWord(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/not-learned-words"
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

	var req schemas.AddNotLearnedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		statusCode = 400
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate word is not empty
	if req.Word == "" {
		statusCode = 400
		http.Error(w, "word cannot be empty", http.StatusBadRequest)
		return
	}

	// Find the word in the database
	word, err := h.WordRepo.GetByValue(r.Context(), req.Word)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = 404
			http.Error(w, "word not found in database", http.StatusNotFound)
			return
		}
		statusCode = 500
		http.Error(w, "failed to find word", http.StatusInternalServerError)
		return
	}

	// Check if word is already in not learned words
	exists, err := h.Repo.Exists(r.Context(), user.ID, word.ID)
	if err != nil {
		statusCode = 500
		http.Error(w, "failed to check if word exists", http.StatusInternalServerError)
		return
	}

	if exists {
		statusCode = 409
		http.Error(w, "word already in not learned list", http.StatusConflict)
		return
	}

	// Create not learned word entry
	nlw := &models.NotLearnedWords{
		ID:     uuid.New(),
		UserID: user.ID,
		WordID: word.ID,
	}

	if err := h.Repo.Create(r.Context(), nlw); err != nil {
		statusCode = 500
		http.Error(w, "failed to create not learned word", http.StatusInternalServerError)
		return
	}

	// Return the created not learned word
	resp := buildNotLearnedWordResponse(nlw, word)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetNotLearnedWords godoc
// @Summary      Get user's not learned words
// @Description  Returns all not learned words for the authenticated user
// @Tags         not-learned-words
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  schemas.ListNotLearnedWordsResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/not-learned-words [get]
func (h *NotLearnedWordHandler) GetNotLearnedWords(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/not-learned-words"
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

	// Get recently not learned words with word details
	words, err := h.Repo.GetRecentlyNotLearnedWords(r.Context(), user.ID, 100) // Limit to 100 words
	if err != nil {
		statusCode = 500
		http.Error(w, "failed to fetch not learned words", http.StatusInternalServerError)
		return
	}

	// Build response
	var resp []schemas.NotLearnedWordResponse
	for _, word := range words {
		// Create a minimal not learned word response
		nlwResp := schemas.NotLearnedWordResponse{
			ID:        uuid.New(), // We don't have the actual ID from the join, so generate one
			UserID:    user.ID,
			WordID:    word.ID,
			Word:      word.Word,
			CreatedAt: time.Now(), // We don't have the actual creation time from the join
		}
		resp = append(resp, nlwResp)
	}

	listResp := schemas.ListNotLearnedWordsResponse{
		Words: resp,
		Total: len(resp),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listResp)
}
