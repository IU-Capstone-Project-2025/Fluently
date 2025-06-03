package handler

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PreferenceHandler struct {
	Repo *postgres.PreferenceRepository
}

func (h PreferenceHandler) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	pref, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "preference not found", http.StatusNotFound)
		return
	}

	resp := schemas.PreferenceResponse{
		ID:             pref.ID,
		CEFRLevel:      pref.CEFRLevel,
		Points:         pref.Points,
		FactEveryday:   pref.FactEveryday,
		Notifications:  pref.Notifications,
		NotificationAt: pref.NotificationsAt,
		WordsPerDay:    pref.WordsPerDay,
		Goal:           pref.Goal,
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *PreferenceHandler) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	var req schemas.CreatePreferenceRequst
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	pref, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "preference not found", http.StatusNotFound)
		return
	}

	pref.CEFRLevel = req.CEFRLevel
	pref.Points = req.Points
	pref.FactEveryday = req.FactEveryday
	pref.Notifications = req.Notifications
	pref.NotificationsAt = req.NotificationAt
	pref.WordsPerDay = req.WordsPerDay
	pref.Goal = req.Goal

	if err := h.Repo.UpdatePreference(r.Context(), pref); err != nil {
		http.Error(w, "failed to update preference", http.StatusInternalServerError)
		return
	}

	resp := schemas.PreferenceResponse{
		ID:             pref.ID,
		CEFRLevel:      pref.CEFRLevel,
		Points:         pref.Points,
		FactEveryday:   pref.FactEveryday,
		Notifications:  pref.Notifications,
		NotificationAt: pref.NotificationsAt,
		WordsPerDay:    pref.WordsPerDay,
		Goal:           pref.Goal,
	}

	json.NewEncoder(w).Encode(resp)
}
