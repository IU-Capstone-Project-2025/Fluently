package schemas

// TelegramLinkRequest is a request body for creating a link token
type TelegramLinkRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required"`
}

// TelegramLinkResponse is a response with a link token
type TelegramLinkResponse struct {
	Token     string `json:"token" example:"abc123xyz"`
	LinkURL   string `json:"link_url" example:"https://example.com/link-google?token=abc123xyz"`
	ExpiresAt string `json:"expires_at" example:"2024-01-01T12:00:00Z"`
}

// TelegramLinkStatusRequest is request body for checking the status of a user by Telegram ID
type TelegramLinkStatusRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required"`
}

// TelegramLinkStatusResponse is a response with a link token
type TelegramLinkStatusResponse struct {
	IsLinked bool       `json:"is_linked"`
	User     *UserBasic `json:"user,omitempty"`
	Message  string     `json:"message"`
}

// UserBasic is a basic user information
type UserBasic struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TelegramUnlinkRequest is a request body for unlinking a Telegram account
type TelegramUnlinkRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required"`
}
