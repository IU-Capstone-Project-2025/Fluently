package postgres

import (
	"context"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenRepository is a repository for refresh tokens
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetByToken finds a refresh token by its token
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

// Revoke marks a refresh token as revoked
func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked", true).Error
}
