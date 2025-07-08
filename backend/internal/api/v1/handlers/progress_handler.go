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

type ProgressHandler struct {
	LearnedWordRepo *postgres.LearnedWordRepository
	WordRepo        *postgres.WordRepository
}

/*
[
  {
    "word": "hello",
    "learned_at": "2024-01-15T10:30:00Z",
    "confidence_score": 85,
    "cnt_reviewed": 3
  },
  {
    "word": "world",
    "learned_at": "2024-01-15T11:45:00Z",
    "confidence_score": 92,
    "cnt_reviewed": 1
  },
  {
    "word": "beautiful",
    "learned_at": "2024-01-16T09:15:00Z",
    "confidence_score": 78,
    "cnt_reviewed": 5
  }
]
*/

type ProgressRequest struct {
	Word            string    `json:"word"`
	LearnedAt       time.Time `json:"learned_at"`
	ConfidenceScore int       `json:"confidence_score"`
	CntReviewed     int       `json:"cnt_reviewed"`
}

// UpdateUserProgress godoc
// @Summary      Update user progress
// @Description  Updates the user's learned words progress. Accepts an array of word progress objects.
// @Tags         progress
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        progress body []ProgressRequest true "List of progress updates"
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
		word, err := h.WordRepo.GetByValue(r.Context(), p.Word)
		if err != nil {
			http.Error(w, "word not found: "+p.Word, http.StatusNotFound)
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

	w.WriteHeader(http.StatusOK)
}
