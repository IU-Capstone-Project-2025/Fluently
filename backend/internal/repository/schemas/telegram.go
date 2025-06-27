package schemas

// TelegramLinkRequest запрос на создание токена связывания
type TelegramLinkRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required"`
}

// TelegramLinkResponse ответ с токеном связывания
type TelegramLinkResponse struct {
	Token     string `json:"token" example:"abc123xyz"`
	LinkURL   string `json:"link_url" example:"https://example.com/link-google?token=abc123xyz"`
	ExpiresAt string `json:"expires_at" example:"2024-01-01T12:00:00Z"`
}

// TelegramLinkStatusRequest запрос проверки статуса пользователя по Telegram ID
type TelegramLinkStatusRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required"`
}

// TelegramLinkStatusResponse ответ о статусе связывания
type TelegramLinkStatusResponse struct {
	IsLinked bool       `json:"is_linked"`
	User     *UserBasic `json:"user,omitempty"`
	Message  string     `json:"message"`
}

// UserBasic базовая информация о пользователе
type UserBasic struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TelegramUnlinkRequest запрос на отвязку Telegram аккаунта
type TelegramUnlinkRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required"`
}
