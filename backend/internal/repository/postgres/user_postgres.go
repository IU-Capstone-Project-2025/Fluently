package postgres

import (
	"context"
	"time"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository is a repository for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID returns a user by id
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail returns a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateRefreshToken updates the refresh token for a user
// This method is deprecated - use RefreshTokenRepository instead
func (r *UserRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	// This method is kept for backward compatibility but should not be used
	// Use RefreshTokenRepository.Create() instead
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", time.Now()).Error
}

// GetByRefreshToken retrieves a user by their refresh token
// This method is deprecated - use RefreshTokenRepository instead
func (r *UserRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error) {
	// This method is kept for backward compatibility but should not be used
	// Use RefreshTokenRepository.GetByToken() instead
	var refreshTokenModel models.RefreshToken
	if err := r.db.WithContext(ctx).First(&refreshTokenModel, "token = ? AND revoked = false AND expires_at > NOW()", refreshToken).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ctx, refreshTokenModel.UserID)
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", time.Now()).Error
}

// ClearRefreshToken clears the refresh token for a user
// This method is deprecated - use RefreshTokenRepository instead
func (r *UserRepository) ClearRefreshToken(ctx context.Context, userID uuid.UUID) error {
	// This method is kept for backward compatibility but should not be used
	// Use RefreshTokenRepository.RevokeByUserID() instead
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// GetByTelegramID finds user by Telegram ID
func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "telegram_id = ?", telegramID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// LinkTelegramID links Telegram ID to existing user
func (r *UserRepository) LinkTelegramID(ctx context.Context, userID uuid.UUID, telegramID int64) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("telegram_id", telegramID).Error
}

// UnlinkTelegramID removes Telegram ID from user
func (r *UserRepository) UnlinkTelegramID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("telegram_id", nil).Error
}
