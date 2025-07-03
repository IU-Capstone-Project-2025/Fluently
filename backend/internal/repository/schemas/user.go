package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email"`
	Provider     string `json:"provider"`
	GoogleID     string `json:"google_id"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	IsActive   bool      `json:"is_active"`
	TelegramID *int64    `json:"telegram_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
