package handlers

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
)

type PreferenceHandler struct {
	Repo *postgres.PreferenceRepository
}

func buildPreferencesResponse(pref *models.Preference) schemas.PreferenceResponse {
	return schemas.PreferenceResponse{
		ID:              pref.ID,
		UserID:          pref.UserID,
		CEFRLevel:       pref.CEFRLevel,
		FactEveryday:    pref.FactEveryday,
		Notifications:   pref.Notifications,
		NotificationsAt: pref.NotificationsAt,
		WordsPerDay:     pref.WordsPerDay,
		Goal:            pref.Goal,
		Subscribed:      pref.Subscribed,
		AvatarImageURL:  pref.AvatarImageURL,
	}
}

func (h *PreferenceHandler) CreateUserPreferences(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var req schemas.CreatePreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pref := &models.Preference{
		ID:              userId,
		UserID:          userId,
		CEFRLevel:       req.CEFRLevel,
		FactEveryday:    req.FactEveryday,
		Notifications:   req.Notifications,
		NotificationsAt: req.NotificationAt,
		WordsPerDay:     req.WordsPerDay,
		Goal:            req.Goal,
		Subscribed:      req.Subscribed,
		AvatarImageURL:  req.AvatarImageURL,
	}

	if err := h.Repo.Create(r.Context(), pref); err != nil {
		http.Error(w, "failed to create preference", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildPreferencesResponse(pref))
}

func (h PreferenceHandler) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	pref, err := h.Repo.GetByID(r.Context(), userId)
	if err != nil {
		http.Error(w, "preference not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPreferencesResponse(pref))
}

func (h *PreferenceHandler) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var req schemas.CreatePreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pref, err := h.Repo.GetByID(r.Context(), userId)
	if err != nil {
		http.Error(w, "preference not found", http.StatusNotFound)
		return
	}

	pref.CEFRLevel = req.CEFRLevel
	pref.FactEveryday = req.FactEveryday
	pref.Notifications = req.Notifications
	pref.NotificationsAt = req.NotificationAt
	pref.WordsPerDay = req.WordsPerDay
	pref.Goal = req.Goal
	pref.Subscribed = req.Subscribed
	pref.AvatarImageURL = req.AvatarImageURL

	if err := h.Repo.Update(r.Context(), pref); err != nil {
		http.Error(w, "failed to update preference", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPreferencesResponse(pref))
}

func (h *PreferenceHandler) DeletePreference(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), userId); err != nil {
		http.Error(w, "failed to delete preference", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
