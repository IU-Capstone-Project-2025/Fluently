package postgres

import (
	"context"
	"time"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LinkTokenRepository is a repository for link tokens
type LinkTokenRepository struct {
	db *gorm.DB
}

// NewLinkTokenRepository creates a new instance of LinkTokenRepository
func NewLinkTokenRepository(db *gorm.DB) *LinkTokenRepository {
	return &LinkTokenRepository{db: db}
}

// Create creates a new link token
func (r *LinkTokenRepository) Create(ctx context.Context, token *models.LinkToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetByToken finds a link token by its token
func (r *LinkTokenRepository) GetByToken(ctx context.Context, token string) (*models.LinkToken, error) {
	var linkToken models.LinkToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&linkToken).Error; err != nil {
		return nil, err
	}
	return &linkToken, nil
}

// MarkAsUsed marks a link token as used
func (r *LinkTokenRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.LinkToken{}).
		Where("id = ?", id).
		Update("used", true).Error
}

// DeleteExpired deletes expired link tokens
func (r *LinkTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.LinkToken{}).Error
}

// GetActiveTelegramTokens returns all active link tokens for a Telegram ID
func (r *LinkTokenRepository) GetActiveTelegramTokens(ctx context.Context, telegramID int64) ([]models.LinkToken, error) {
	var tokens []models.LinkToken
	err := r.db.WithContext(ctx).
		Where("telegram_id = ? AND used = false AND expires_at > ?", telegramID, time.Now()).
		Find(&tokens).Error
	return tokens, err
}
