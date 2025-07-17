package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
)

// swagger:ignore
var _ schemas.ErrorResponse

// ------------------- request / response schemas --------------------

type ChatMessage struct {
	Author  string `json:"author"` // "user" or "llm"
	Message string `json:"message"`
}

type ChatRequest struct {
	Chat []ChatMessage `json:"chat"`
}

type ChatResponse struct {
	Chat []ChatMessage `json:"chat"`
}

// -------------------- handler --------------------------------------

type ChatHandler struct {
	Redis       *goredis.Client
	HistoryRepo *postgres.ChatHistoryRepository
	LLMClient   *utils.LLMClient
}

var stopWords = []string{"finish", "всё", "хочу закончить"}

// Chat godoc
// @Summary Отправить сообщение в диалоге с ИИ
// @Description Добавляет очередное сообщение пользователя, получает ответ LLM, сохраняет историю в Redis. При обнаружении стоп-слов диалог считается завершённым и будет сохранён в базу.
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body ChatRequest true "Сообщения диалога"
// @Success 200 {object} ChatResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /api/v1/chat [post]
func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if len(req.Chat) == 0 {
		http.Error(w, "chat array empty", http.StatusBadRequest)
		return
	}

	user, err := utils.GetCurrentUser(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert messages to LLM format
	var llmMsgs []utils.LLMMessage
	for _, m := range req.Chat {
		role := "user"
		if m.Author == "llm" {
			role = "assistant"
		}
		if m.Author == "system" { // in case
			role = "system"
		}
		llmMsgs = append(llmMsgs, utils.LLMMessage{Role: role, Content: m.Message})
	}

	// Call LLM
	reply, err := h.LLMClient.Chat(ctx, llmMsgs, "balanced", nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append reply to chat
	req.Chat = append(req.Chat, ChatMessage{Author: "llm", Message: reply})

	// Store to redis
	key := "chat:" + user.ID.String()
	if data, _ := json.Marshal(req.Chat); data != nil {
		h.Redis.Set(ctx, key, data, 24*time.Hour) // expire after a day
	}

	// Check stop words on the last user message
	lastMsg := strings.ToLower(req.Chat[len(req.Chat)-2].Message)
	shouldFinish := false
	for _, sw := range stopWords {
		if strings.Contains(lastMsg, sw) {
			shouldFinish = true
			break
		}
	}

	if shouldFinish {
		h.finishChat(ctx, user.ID)
	}

	// Return chat with AI reply
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Chat: req.Chat})
}

// FinishChat godoc
// @Summary Завершить диалог с ИИ
// @Description Принудительно завершает текущий диалог: переносит историю из Redis в Postgres и очищает кеш.
// @Tags Chat
// @Produce json
// @Success 204 "Диалог сохранён, ответ без тела"
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /api/v1/chat/finish [post]
func (h *ChatHandler) FinishChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := utils.GetCurrentUser(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.finishChat(ctx, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// finishChat flushes redis to postgres and clears key
func (h *ChatHandler) finishChat(ctx context.Context, userID uuid.UUID) error {
	key := "chat:" + userID.String()
	data, err := h.Redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil // nothing to flush
	}

	history := models.ChatHistory{
		UserID:   userID,
		Messages: datatypes.JSON(data),
	}
	if err := h.HistoryRepo.Create(ctx, &history); err != nil {
		return err
	}

	h.Redis.Del(ctx, key)
	return nil
}
