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
)

// PreferenceHandler is a handler for preferences
type PreferenceHandler struct {
	Repo *postgres.PreferenceRepository
}

// buildPreferencesResponse builds a response from a preference
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

// CreateUserPreferences godoc
func (h *PreferenceHandler) CreateUserPreferences(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/preferences"
	method := r.Method
	statusCode := 201
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()

	userId, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		statusCode = 400
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var req schemas.CreatePreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		statusCode = 400
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
		statusCode = 500
		http.Error(w, "failed to create preference", http.StatusInternalServerError)
		return
	}

	// Return the created preference
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildPreferencesResponse(pref))
}

// GetUserPreferences godoc
// @Summary Get user preferences
// @Description Retrieves user preferences
// @Tags preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} schemas.PreferenceResponse "Successfully retrieved preferences"
// @Failure 400 {string} string "Bad request - invalid user or preferences"
// @Failure 401 {string} string "Unauthorized - invalid or missing token"
// @Failure 404 {string} string "Not found - user or preferences not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/preferences [get]
func (h PreferenceHandler) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/preferences"
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

	pref, err := h.Repo.GetByUserID(r.Context(), user.ID)
	if err != nil {
		statusCode = 404
		http.Error(w, "preference not found", http.StatusNotFound)
		return
	}

	// Return the preferences
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPreferencesResponse(pref))
}

// UpdateUserPreferences godoc
// @Summary Update user preferences
// @Description Updates user preferences
// @Tags preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} schemas.PreferenceResponse "Successfully updated preferences"
// @Failure 400 {string} string "Bad request - invalid user or preferences"
// @Failure 401 {string} string "Unauthorized - invalid or missing token"
// @Failure 404 {string} string "Not found - user or preferences not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/preferences [put]
func (h *PreferenceHandler) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/preferences"
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

	var req schemas.UpdatePreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		statusCode = 400
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pref, err := h.Repo.GetByUserID(r.Context(), user.ID)
	if err != nil {
		statusCode = 404
		http.Error(w, "preference not found", http.StatusNotFound)
		return
	}

	// ======================== Update block ========================
	if req.CEFRLevel != nil {
		pref.CEFRLevel = *req.CEFRLevel
	}
	if req.FactEveryday != nil {
		pref.FactEveryday = *req.FactEveryday
	}
	if req.Notifications != nil {
		pref.Notifications = *req.Notifications
	}
	if req.NotificationAt != nil {
		pref.NotificationsAt = req.NotificationAt
	}
	if req.WordsPerDay != nil {
		pref.WordsPerDay = *req.WordsPerDay
	}
	if req.Goal != nil {
		pref.Goal = *req.Goal
	}
	if req.Subscribed != nil {
		pref.Subscribed = *req.Subscribed
	}
	if req.AvatarImageURL != nil {
		pref.AvatarImageURL = *req.AvatarImageURL
	}
	// ============================================================

	if err := h.Repo.Update(r.Context(), pref.ID, &req); err != nil {
		statusCode = 500
		http.Error(w, "failed to update preference", http.StatusInternalServerError)
		return
	}

	// Return the updated preferences
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPreferencesResponse(pref))
}

// DeleteUserPreferences deletes user preferences
func (h *PreferenceHandler) DeletePreference(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/preferences"
	method := r.Method
	statusCode := 204
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()

	userId, err := utils.ParseUUIDParam(r, "user_id")
	if err != nil {
		statusCode = 400
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), userId); err != nil {
		statusCode = 500
		http.Error(w, "failed to delete preference", http.StatusInternalServerError)
		return
	}

	// Return no content
	w.WriteHeader(http.StatusNoContent)
}
