package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	schemas "fluently/go-backend/internal/repository/schemas"

	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"
)

var _ = schemas.ErrorResponse{}

// ChatHistoryItem is used in API response
// swagger:model
type ChatHistoryItem struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	Chat      []ChatMessage `json:"chat"`
}

type ChatHistoryHandler struct {
	Repo *postgres.ChatHistoryRepository
}

// GetHistory godoc
// @Summary Получить историю диалогов за выбранный день
// @Description Возвращает все завершённые диалоги пользователя за указанный день (UTC). Передайте ?day=YYYY-MM-DDTHH:MM:SSZ
// @Tags Chat
// @Produce json
// @Param day query string true "Точка внутри нужного дня в формате RFC3339 (пример: 2025-07-17T00:00:00Z)"
// @Success 200 {array} ChatHistoryItem
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Router /api/v1/chat/history [get]
func (h *ChatHistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := utils.GetCurrentUser(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dayStr := r.URL.Query().Get("day")
	if dayStr == "" {
		http.Error(w, "query param 'day' required", http.StatusBadRequest)
		return
	}
	dayTime, err := time.Parse(time.RFC3339, dayStr)
	if err != nil {
		http.Error(w, "invalid day format", http.StatusBadRequest)
		return
	}

	dayStart := time.Date(dayTime.Year(), dayTime.Month(), dayTime.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	histories, err := h.Repo.ListByUserAndDay(ctx, user.ID, dayStart, dayEnd)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	var resp []ChatHistoryItem
	for _, hst := range histories {
		var chat []ChatMessage
		_ = json.Unmarshal(hst.Messages, &chat)
		resp = append(resp, ChatHistoryItem{
			ID:        hst.ID.String(),
			CreatedAt: hst.CreatedAt,
			Chat:      chat,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
