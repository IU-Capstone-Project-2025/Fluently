package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"
	"fluently/go-backend/pkg/logger"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"

	"github.com/bsm/redislock"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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

	// Acquire per-user distributed lock so concurrent requests (e.g. /chat and
	// /chat/finish) for the same user are serialized.
	user, err := utils.GetCurrentUser(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	lock, err := utils.AcquireChatLock(ctx, user.ID)
	if err == redislock.ErrNotObtained {
		http.Error(w, "another chat operation is in progress", http.StatusTooManyRequests)
		return
	} else if err != nil {
		logger.Log.Error("failed to acquire chat lock", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer lock.Release(ctx)

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
	logger.Log.Info("LLM reply", zap.String("reply", reply))
	if reply == "" {
		logger.Log.Error("LLM reply is empty")
		http.Error(w, "LLM reply is empty", http.StatusInternalServerError)
		return
	}
	if err != nil {
		logger.Log.Error("LLM error", zap.Error(err))
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
	// Check out OutOfRange error from the LLM
	lastMsg := ""
	if len(req.Chat) > 1 {
		lastMsg = strings.ToLower(req.Chat[len(req.Chat)-2].Message)
	} else {
		logger.Log.Error("no last message")
		http.Error(w, "no last message", http.StatusInternalServerError)
		return
	}
	shouldFinish := false
	for _, sw := range stopWords {
		if strings.Contains(lastMsg, sw) {
			shouldFinish = true
			break
		}
	}

	if shouldFinish {
		if err := h.flushChat(ctx, user.ID); err != nil {
			logger.Log.Error("failed to flush chat", zap.Error(err))
		}
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

	lock, err := utils.AcquireChatLock(ctx, user.ID)
	if err == redislock.ErrNotObtained {
		http.Error(w, "chat operation in progress, try later", http.StatusTooManyRequests)
		return
	} else if err != nil {
		logger.Log.Error("failed to acquire chat lock", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer lock.Release(ctx)

	if err := h.flushChat(ctx, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// flushChat moves chat history from Redis to Postgres atomically. Callers must
// hold the per-user chat lock before invoking this method.
func (h *ChatHandler) flushChat(ctx context.Context, userID uuid.UUID) error {
	key := "chat:" + userID.String()

	// Atomically fetch and delete using GETDEL (Redis >= 6.2).
	data, err := h.Redis.GetDel(ctx, key).Bytes()
	if err == goredis.Nil {
		return nil // nothing to flush
	}
	if err != nil {
		logger.Log.Error("failed to get/del chat from redis", zap.Error(err))
		return err
	}
	if len(data) == 0 {
		return nil
	}

	history := models.ChatHistory{
		UserID:   userID,
		Messages: datatypes.JSON(data),
	}
	return h.HistoryRepo.Create(ctx, &history)
}
