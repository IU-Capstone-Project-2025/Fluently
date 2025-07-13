package handlers

import (
	"encoding/json"
	"net/http"

	_ "fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"
)

type DayWordHandler struct {
	WordRepo       *postgres.WordRepository
	PreferenceRepo *postgres.PreferenceRepository
}

type DayWordResponse struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

// godoc
// @Summary      Get day word
// @Description  Returns the day word for the user
// @Tags         day-word
// @Produce      json
// @Security     BearerAuth
// @Success 	 200 {object}  DayWordResponse "Successfully returned day word"
// @Failure      400  {string}  string  "Invalid request - plain text error message"
// @Failure      404  {string}  string  "Resource not found - plain text error message"
// @Failure      500  {string}  string  "Internal server error - plain text error message"
// @Router       /api/v1/day-word/ [get]
func (h *DayWordHandler) GetDayWord(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := user.ID

	userPref, err := h.PreferenceRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get preference", http.StatusInternalServerError)
		return
	}

	dayWord, err := h.WordRepo.GetDayWord(r.Context(), userPref.CEFRLevel, userID)
	if err != nil {
		http.Error(w, "failed to get day word", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DayWordResponse{
		Word:        dayWord.Word,
		Translation: dayWord.Translation,
	})
}
