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
		ID:             pref.ID,
		CEFRLevel:      pref.CEFRLevel,
		FactEveryday:   pref.FactEveryday,
		Notifications:  pref.Notifications,
		NotificationAt: pref.NotificationsAt,
		WordsPerDay:    pref.WordsPerDay,
		Goal:           pref.Goal,
		Subscribed:     pref.Subscribed,
		AvatarImage:    pref.AvatarImage,
	}
}

// CreateUserPreferences godoc
// @Summary      Создать предпочтения пользователя
// @Description  Создаёт предпочтения пользователя по его ID
// @Tags         preferences
// @Accept       json
// @Produce      json
// @Param        id         path      string                          true  "ID пользователя"
// @Param        preference body      schemas.CreatePreferenceRequst  true  "Данные предпочтений"
// @Success      201  {object}  schemas.PreferenceResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{id}/preferences/ [post]
func (h *PreferenceHandler) CreateUserPreferences(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	var req schemas.CreatePreferenceRequst
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	pref := &models.Preference{
		ID:              id,
		CEFRLevel:       req.CEFRLevel,
		FactEveryday:    req.FactEveryday,
		Notifications:   req.Notifications,
		NotificationsAt: req.NotificationAt,
		WordsPerDay:     req.WordsPerDay,
		Goal:            req.Goal,
		Subscribed:      req.Subscribed,
		AvatarImage:     req.AvatarImage,
	}

	if err := h.Repo.Create(r.Context(), pref); err != nil {
		http.Error(w, "failed to create preference", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildPreferencesResponse(pref))
}

// GetUserPreferences godoc
// @Summary      Получить предпочтения пользователя
// @Description  Возвращает предпочтения пользователя по его ID
// @Tags         preferences
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Success      200  {object}  schemas.PreferenceResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Router       /users/{id}/preferences/ [get]
func (h PreferenceHandler) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	pref, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "preference not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(pref)
}

// UpdateUserPreferences godoc
// @Summary      Обновить предпочтения пользователя
// @Description  Обновляет предпочтения пользователя по его ID
// @Tags         preferences
// @Accept       json
// @Produce      json
// @Param        id         path      string                          true  "ID пользователя"
// @Param        preference body      schemas.CreatePreferenceRequst  true  "Данные предпочтений"
// @Success      200  {object}  schemas.PreferenceResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{id}/preferences/ [put]
func (h *PreferenceHandler) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
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
	pref.FactEveryday = req.FactEveryday
	pref.Notifications = req.Notifications
	pref.NotificationsAt = req.NotificationAt
	pref.WordsPerDay = req.WordsPerDay
	pref.Goal = req.Goal
	pref.Subscribed = req.Subscribed
	pref.AvatarImage = req.AvatarImage

	if err := h.Repo.Update(r.Context(), pref); err != nil {
		http.Error(w, "failed to update preference", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pref)
}
