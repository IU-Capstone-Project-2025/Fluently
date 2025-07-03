package postgres

import (
	"context"
	"testing"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

	// Refresh Token
	newToken := "new_refresh_token_test"
	err = userRepo.UpdateRefreshToken(ctx, user.ID, newToken)
	assert.NoError(t, err)

	// Get User by Token
	found, err := userRepo.GetByRefreshToken(ctx, newToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, newToken, found.RefreshToken)
	assert.False(t, found.LastLoginAt.IsZero()) // Check if update
}

func TestClearRefreshToken(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:           uuid.New(),
		Name:         "Clear User",
		Email:        "clear@example.com",
		RefreshToken: "token_to_clear",
		IsActive:     true,
	}

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	err = userRepo.ClearRefreshToken(ctx, user.ID)
	assert.NoError(t, err)

	found, err := userRepo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "", found.RefreshToken)
}

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
