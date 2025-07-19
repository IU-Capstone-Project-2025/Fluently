package postgres

import (
	"context"
	"testing"
	"time"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCreateAndGetUser
func TestCreateAndGetUser(t *testing.T) {
	ctx := context.Background()

	id := uuid.New()
	user := &models.User{
		ID:       id,
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     "user",
		IsActive: true,
	}

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	found, err := userRepo.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Name, found.Name)
}

// TestUpdateRefreshToken
func TestUpdateRefreshToken(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:       uuid.New(),
		Name:     "Refresh User",
		Email:    "refresh@example.com",
		Role:     "user",
		IsActive: true,
	}
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	// Create refresh token using the new repository
	newToken := "new_refresh_token_test"
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = refreshTokenRepo.Create(ctx, refreshToken)
	assert.NoError(t, err)

	// Get User by Token (using the deprecated but working method)
	found, err := userRepo.GetByRefreshToken(ctx, newToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.False(t, found.LastLoginAt.IsZero()) // Check if last login was updated
}

// TestClearRefreshToken clears the refresh token
func TestClearRefreshToken(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:       uuid.New(),
		Name:     "Clear User",
		Email:    "clear@example.com",
		IsActive: true,
	}

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	// Create a refresh token first
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     "token_to_clear",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = refreshTokenRepo.Create(ctx, refreshToken)
	assert.NoError(t, err)

	// Clear the refresh token
	err = userRepo.ClearRefreshToken(ctx, user.ID)
	assert.NoError(t, err)

	// Verify the token was revoked
	foundToken, err := refreshTokenRepo.GetByToken(ctx, "token_to_clear")
	assert.Error(t, err) // Should not find the token as it's revoked
	assert.Nil(t, foundToken)
}

// TestLinkAndGetByTelegramID is a test for linking a user to a Telegram ID
func TestLinkAndGetByTelegramID(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:       uuid.New(),
		Name:     "Telegram User",
		Email:    "telegram@example.com",
		IsActive: true,
	}
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	telegramID := int64(1234567890)
	err = userRepo.LinkTelegramID(ctx, user.ID, telegramID)
	assert.NoError(t, err)

	found, err := userRepo.GetByTelegramID(ctx, telegramID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)

	// Unlink
	err = userRepo.UnlinkTelegramID(ctx, user.ID)
	assert.NoError(t, err)

	// Check if telegramID became nil
	found, err = userRepo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Nil(t, found.TelegramID)
}

// TestDuplicateEmail is a test for checking if a user with the same email already exists
func TestDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	email := "duplicate@example.com"

	first := &models.User{
		ID:       uuid.New(),
		Name:     "Original User",
		Email:    email,
		IsActive: true,
	}
	err := userRepo.Create(ctx, first)
	assert.NoError(t, err)

	second := &models.User{
		ID:       uuid.New(),
		Name:     "Duplicate User",
		Email:    email,
		IsActive: true,
	}
	err = userRepo.Create(ctx, second)
	assert.Error(t, err, "expected unique constraint violation")
}

// TestRefreshTokenRepository tests the new refresh token repository
func TestRefreshTokenRepository(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:       uuid.New(),
		Name:     "Token Test User",
		Email:    "token@example.com",
		Role:     "user",
		IsActive: true,
	}
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	// Create refresh token
	token := "test_refresh_token_123"
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = refreshTokenRepo.Create(ctx, refreshToken)
	assert.NoError(t, err)

	// Get token by token value
	foundToken, err := refreshTokenRepo.GetByToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundToken.UserID)
	assert.Equal(t, token, foundToken.Token)
	assert.False(t, foundToken.Revoked)

	// Get tokens by user ID
	userTokens, err := refreshTokenRepo.GetByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Len(t, userTokens, 1)
	assert.Equal(t, token, userTokens[0].Token)

	// Revoke token
	err = refreshTokenRepo.RevokeToken(ctx, token)
	assert.NoError(t, err)

	// Verify token is revoked
	foundToken, err = refreshTokenRepo.GetByToken(ctx, token)
	assert.Error(t, err) // Should not find revoked token
	assert.Nil(t, foundToken)

	// Create another token and revoke by user ID
	token2 := "test_refresh_token_456"
	refreshToken2 := &models.RefreshToken{
		UserID:    user.ID,
		Token:     token2,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = refreshTokenRepo.Create(ctx, refreshToken2)
	assert.NoError(t, err)

	err = refreshTokenRepo.RevokeByUserID(ctx, user.ID)
	assert.NoError(t, err)

	// Verify all user tokens are revoked
	userTokens, err = refreshTokenRepo.GetByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Len(t, userTokens, 0) // No active tokens
}
