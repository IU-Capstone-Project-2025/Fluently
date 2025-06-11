package schemas

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name     string    `json:"name" binding:"required"`
	SubLevel bool      `json:"sub_level"`
	PrefID   uuid.UUID `json:"pref_id"`
}

type UserResponse struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	SubLevel bool            `json:"sub_level"`
	Pref     *PreferenceMini `json:"preference,omitempty"`
}

type PreferenceMini struct {
	ID        uuid.UUID `json:"id"`
	CEFRLevel string    `json:"ceft_level"`
}
