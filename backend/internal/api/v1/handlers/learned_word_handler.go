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

type LearnedWordHandler struct {
	Repo *postgres.LearnedWordRepository
}

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

// GetLearnedWords godoc
// @Summary      Get user's learned words
// @Description  Returns all learned words for a user
// @Tags         learned-words
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id   path      string  true  "User ID"
// @Success      200  {array}   schemas.LearnedWordResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{user_id}/learned-words/ [get]
func (h *LearnedWordHandler) GetLearnedWords(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	words, err := h.Repo.ListByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch learned words", http.StatusInternalServerError)
		return
	}

	var resp []schemas.LearnedWordResponse
	for _, w := range words {
		resp = append(resp, buildLearnedWordResponse(&w))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetLearnedWord godoc
// @Summary      Get user's learned word by word ID
// @Description  Returns a learned word by user ID and word ID
// @Tags         learned-words
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id   path      string  true  "User ID"
// @Param        word_id   path      string  true  "Word ID"
// @Success      200  {object}  schemas.LearnedWordResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Router       /users/{user_id}/learned-words/{word_id} [get]
func (h *LearnedWordHandler) GetLearnedWord(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	wordID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	word, err := h.Repo.GetByUserWordID(r.Context(), userID, wordID)
	if err != nil {
		http.Error(w, "learned word not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildLearnedWordResponse(word))
}

// CreateLearnedWord godoc
// @Summary      Add a learned word
// @Description  Adds a new learned word for a user
// @Tags         learned-words
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        learned_word  body      schemas.CreateLearnedWordRequest  true  "Learned word data"
// @Success      201  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{user_id}/learned-words/ [post]
func (h *LearnedWordHandler) CreateLearnedWord(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var req schemas.CreateLearnedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID != uuid.Nil && req.UserID != userID {
		http.Error(w, "user_id mismatch", http.StatusBadRequest)
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
		http.Error(w, "failed to create learned word", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// UpdateLearnedWord godoc
// @Summary      Update a learned word
// @Description  Updates a user's learned word by word ID
// @Tags         learned-words
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id       path      string                        true  "User ID"
// @Param        word_id       path      string                        true  "Word ID"
// @Param        learned_word  body      schemas.CreateLearnedWordRequest true  "Learned word data"
// @Success      200  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{user_id}/learned-words/{word_id} [put]
func (h *LearnedWordHandler) UpdateLearnedWord(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	wordID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	var req schemas.CreateLearnedWordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	word, err := h.Repo.GetByUserWordID(r.Context(), userID, wordID)
	if err != nil {
		http.Error(w, "learned word not found", http.StatusNotFound)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteLearnedWord godoc
// @Summary      Delete a learned word
// @Description  Deletes a user's learned word by word ID
// @Tags         learned-words
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id   path      string  true  "User ID"
// @Param        word_id   path      string  true  "Word ID"
// @Success      204  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{user_id}/learned-words/{word_id} [delete]
func (h *LearnedWordHandler) DeleteLearnedWord(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	wordID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), userID, wordID); err != nil {
		http.Error(w, "failed to delete learned word", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
