package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreatePreferenceRequest struct {
	UserID         uuid.UUID  `json:"user_id"`
	CEFRLevel      string     `json:"cefr_level" binding:"required"`
	FactEveryday   bool       `json:"fact_everyday"`
	Notifications  bool       `json:"notifications"`
	NotificationAt *time.Time `json:"notification_at,omitempty"`
	WordsPerDay    int        `json:"words_per_day"`
	Goal           string     `json:"goal"`
	Subscribed     bool       `json:"subscribed"`
	AvatarImageURL string     `json:"avatar_image_url"`
}

type UpdatePreferenceRequest struct {
	CEFRLevel      *string    `json:"cefr_level,omitempty"`
	FactEveryday   *bool      `json:"fact_everyday,omitempty"`
	Notifications  *bool      `json:"notifications,omitempty"`
	NotificationAt *time.Time `json:"notification_at,omitempty"`
	WordsPerDay    *int       `json:"words_per_day,omitempty"`
	Goal           *string    `json:"goal,omitempty"`
	Subscribed     *bool      `json:"subscribed,omitempty"`
	AvatarImageURL *string    `json:"avatar_image_url,omitempty"`
}

type PreferenceResponse struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	CEFRLevel       string     `json:"cefr_level"`
	FactEveryday    bool       `json:"fact_everyday"`
	Notifications   bool       `json:"notifications"`
	NotificationsAt *time.Time `json:"notification_at,omitempty"`
	WordsPerDay     int        `json:"words_per_day"`
	Goal            string     `json:"goal"`
	Subscribed      bool       `json:"subscribed"`
	AvatarImageURL  string     `json:"avatar_image_url"`
}
