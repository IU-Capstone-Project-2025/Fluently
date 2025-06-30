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
	NotificationAt *time.Time `json:"notification_at"`
	WordsPerDay    int        `json:"words_per_day"`
	Goal           string     `json:"goal"`
	Subscribed     bool       `json:"subscribed"`
	AvatarImageURL string     `json:"avatar_image"`
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
	AvatarImageURL  string     `json:"avatar_image,omitempty"`
}
