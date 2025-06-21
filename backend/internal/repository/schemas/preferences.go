package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreatePreferenceRequest struct {
	CEFRLevel      float64    `json:"cefr_level" binding:"required"`
	FactEveryday   bool       `json:"fact_everyday"`
	Notifications  bool       `json:"notifications"`
	NotificationAt *time.Time `json:"notification_at"`
	WordsPerDay    int        `json:"words_per_day"`
	Goal           string     `json:"goal"`
	Subscribed     bool       `json:"subscribed"`
	AvatarImage    []byte     `json:"avata_image"` // Base64
}

type PreferenceResponse struct {
	ID              uuid.UUID  `json:"id"`
	CEFRLevel       float64    `json:"cefr_level"`
	FactEveryday    bool       `json:"fact_everyday"`
	Notifications   bool       `json:"notifications"`
	NotificationsAt *time.Time `json:"notification_at,omitempty"`
	WordsPerDay     int        `json:"words_per_day"`
	Goal            string     `json:"goal"`
	Subscribed      bool       `json:"subscribed"`
	AvatarImage     []byte     `json:"avatar_image,omitempty"`
}
