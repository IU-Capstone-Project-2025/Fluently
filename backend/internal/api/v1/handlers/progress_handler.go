package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProgressHandler handles the progress endpoint
type ProgressHandler struct {
	LearnedWordRepo *postgres.LearnedWordRepository
	WordRepo        *postgres.WordRepository
}

// ProgressRequest is a request body for updating user progress
type ProgressRequest struct {
	WordID          uuid.UUID `json:"word_id"`
	LearnedAt       time.Time `json:"learned_at"`
	ConfidenceScore int       `json:"confidence_score"`
	CntReviewed     int       `json:"cnt_reviewed"`
}

// UpdateUserProgress godoc
// @Summary      Update user progress
// @Description  Updates the user's learned words progress. Accepts an array of word progress objects with word-translation pairs.
// @Tags         progress
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        progress body []ProgressRequest true "List of progress updates with word-translation pairs"
// @Success      200  {string}  string  "ok"
// @Router       /api/v1/progress [post]
func (h *ProgressHandler) UpdateUserProgress(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := user.ID

	var progress []ProgressRequest
	err = json.NewDecoder(r.Body).Decode(&progress)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for _, p := range progress {
		word, err := h.WordRepo.GetByID(r.Context(), p.WordID)
		if err != nil {
			http.Error(w, "word not found: "+err.Error(), http.StatusNotFound)
			return
		}

		existing, err := h.LearnedWordRepo.GetByUserWordID(r.Context(), userID, word.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			http.Error(w, "failed to get learned word", http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()

		if existing == nil {
			lw := &models.LearnedWords{
				ID:               uuid.New(),
				UserID:           userID,
				WordID:           word.ID,
				LearnedAt:        p.LearnedAt,
				LastReviewed:     now,
				CountOfRevisions: p.CntReviewed,
				ConfidenceScore:  p.ConfidenceScore,
			}
			if err := h.LearnedWordRepo.Create(r.Context(), lw); err != nil {
				http.Error(w, "failed to create learned word", http.StatusInternalServerError)
				return
			}
		} else {
			existing.LearnedAt = p.LearnedAt
			existing.LastReviewed = now
			existing.CountOfRevisions = p.CntReviewed
			existing.ConfidenceScore = p.ConfidenceScore

			if err := h.LearnedWordRepo.Update(r.Context(), existing); err != nil {
				http.Error(w, "failed to update learned word", http.StatusInternalServerError)
				return
			}
		}
	}

	// Return ok
	w.WriteHeader(http.StatusOK)
}
