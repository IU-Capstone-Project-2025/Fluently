package postgres

import (
	"context"
	"encoding/json"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/pkg/logger"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ChatHistoryRepository handles chat_history table operations.
type ChatHistoryRepository struct {
	db *gorm.DB
}

func NewChatHistoryRepository(db *gorm.DB) *ChatHistoryRepository {
	return &ChatHistoryRepository{db: db}
}

// Create inserts a finished chat history.
func (r *ChatHistoryRepository) Create(ctx context.Context, history *models.ChatHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// // ListByUser fetches histories for a user.
// func (r *ChatHistoryRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.ChatHistory, error) {
// 	var list []models.ChatHistory
// 	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")
// 	if limit > 0 {
// 		query = query.Limit(limit)
// 	}
// 	err := query.Find(&list).Error
// 	return list, err
// }

// ListByUserAndDay returns histories for a specific user created on a given UTC day.
func (r *ChatHistoryRepository) ListByUserAndDay(ctx context.Context, userID uuid.UUID, dayStart, dayEnd time.Time) ([]models.ChatHistory, error) {
	var list []models.ChatHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, dayStart, dayEnd).
		Order("created_at ASC").
		Find(&list).Error
	return list, err
}

// ToJSON helper to marshal messages directly.
func ToJSON(v any) datatypes.JSON {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Log.Error("failed to marshal to json", zap.Error(err))
		return datatypes.JSON{}
	}
	return datatypes.JSON(b)
}
