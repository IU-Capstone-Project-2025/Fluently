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
func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(refreshToken).Error
}

// GetByToken retrieves a refresh token by its token value
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.WithContext(ctx).First(&refreshToken, "token = ? AND revoked = false AND expires_at > NOW()", token).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// GetByUserID retrieves all active refresh tokens for a user
func (r *RefreshTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.RefreshToken, error) {
	var refreshTokens []models.RefreshToken
	if err := r.db.WithContext(ctx).Where("user_id = ? AND revoked = false AND expires_at > NOW()", userID).Find(&refreshTokens).Error; err != nil {
		return nil, err
	}
	return refreshTokens, nil
}

// RevokeToken revokes a specific refresh token
func (r *RefreshTokenRepository) RevokeToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

// RevokeByUserID revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

// DeleteExpired removes expired refresh tokens
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < NOW()").Delete(&models.RefreshToken{}).Error
}

// Update updates a refresh token
func (r *RefreshTokenRepository) Update(ctx context.Context, refreshToken *models.RefreshToken) error {
	return r.db.WithContext(ctx).Save(refreshToken).Error
}

// Delete deletes a refresh token
func (r *RefreshTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.RefreshToken{}, "id = ?", id).Error
}
