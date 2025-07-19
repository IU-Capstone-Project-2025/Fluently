package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProgressHandler handles the progress endpoint
type ProgressHandler struct {
	LearnedWordRepo    *postgres.LearnedWordRepository
	WordRepo           *postgres.WordRepository
	NotLearnedWordRepo *postgres.NotLearnedWordRepository
}

// ProgressRequest is a request body for updating user progress
type ProgressRequest struct {
	WordID          uuid.UUID  `json:"word_id"`
	LearnedAt       *time.Time `json:"learned_at,omitempty"`
	ConfidenceScore *int       `json:"confidence_score,omitempty"`
	CntReviewed     *int       `json:"cnt_reviewed,omitempty"`
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
	start := time.Now()
	endpoint := "/api/v1/progress"
	method := r.Method
	statusCode := 200
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()

	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		statusCode = 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := user.ID

	var progress []ProgressRequest
	err = json.NewDecoder(r.Body).Decode(&progress)
	if err != nil {
		statusCode = 400
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(progress) == 0 {
		statusCode = 400
		http.Error(w, "no progress provided", http.StatusBadRequest)
		return
	}

	for _, p := range progress {
		word, err := h.WordRepo.GetByID(r.Context(), p.WordID)
		if err != nil {
			statusCode = 404
			http.Error(w, "word not found: "+err.Error(), http.StatusNotFound)
			return
		}

		existing, err := h.LearnedWordRepo.GetByUserWordID(r.Context(), userID, word.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			statusCode = 500
			http.Error(w, "failed to get learned word", http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()

		if isNotLearnedWord(p) {
			notLearnedExists, err := h.NotLearnedWordRepo.Exists(r.Context(), userID, word.ID)
			if err != nil {
				statusCode = 500
				http.Error(w, "failed to check if not learned word exists", http.StatusInternalServerError)
				return
			}

			if !notLearnedExists {
				nlw := &models.NotLearnedWords{
					ID:     uuid.New(),
					UserID: userID,
					WordID: word.ID,
				}
				if err := h.NotLearnedWordRepo.Create(r.Context(), nlw); err != nil {
					statusCode = 500
					http.Error(w, "failed to create not learned word", http.StatusInternalServerError)
					return
				}
			}

			continue
		}

		if existing == nil {
			lw := &models.LearnedWords{
				ID:           uuid.New(),
				UserID:       userID,
				WordID:       word.ID,
				LastReviewed: now,
			}

			if p.LearnedAt != nil {
				lw.LearnedAt = *p.LearnedAt
			} else {
				lw.LearnedAt = now
			}

			if p.ConfidenceScore != nil {
				lw.ConfidenceScore = *p.ConfidenceScore
			} else {
				lw.ConfidenceScore = 0
			}

			if p.CntReviewed != nil {
				lw.CountOfRevisions = *p.CntReviewed
			} else {
				lw.CountOfRevisions = 0
			}

			if err := h.LearnedWordRepo.Create(r.Context(), lw); err != nil {
				statusCode = 500
				http.Error(w, "failed to create learned word", http.StatusInternalServerError)
				return
			}

			if err := h.NotLearnedWordRepo.DeleteIfExists(r.Context(), userID, word.ID); err != nil {
				statusCode = 500
				http.Error(w, "failed to delete not learned word", http.StatusInternalServerError)
				return
			}
		} else {
			existing.LastReviewed = now

			if p.LearnedAt != nil {
				existing.LearnedAt = *p.LearnedAt
			}

			if p.ConfidenceScore != nil {
				existing.ConfidenceScore = *p.ConfidenceScore
			}

			if p.CntReviewed != nil {
				existing.CountOfRevisions = *p.CntReviewed
			}

			if err := h.LearnedWordRepo.Update(r.Context(), existing); err != nil {
				statusCode = 500
				http.Error(w, "failed to update learned word", http.StatusInternalServerError)
				return
			}
		}
	}

	// Return ok
	w.WriteHeader(http.StatusOK)
}

func isNotLearnedWord(p ProgressRequest) bool {
	return p.LearnedAt == nil && p.ConfidenceScore == nil && p.CntReviewed == nil
}
