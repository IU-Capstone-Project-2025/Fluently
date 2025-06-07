package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreatePreferenceRequst struct {
	CEFRLevel      string     `json:"cefr_level" binding:"required"`
	Points         int        `json:"points"`
	FactEveryday   bool       `json:"fact_everyday"`
	Notifications  bool       `json:"notifications"`
	NotificationAt *time.Time `json:"notification_at"`
	WordsPerDay    int        `json:"words_per_day"`
	Goal           string     `json:"goal"`
}

type PreferenceResponse struct {
	ID             uuid.UUID  `json:"id"`
	CEFRLevel      string     `json:"cefr_level"`
	Points         int        `json:"points"`
	FactEveryday   bool       `json:"fact_everyday"`
	Notifications  bool       `json:"notifications"`
	NotificationAt *time.Time `json:"notification_at,omitempty"`
	WordsPerDay    int        `json:"words_per_day"`
	Goal           string     `json:"goal"`
}
