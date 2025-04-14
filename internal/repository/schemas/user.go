package schemas

<<<<<<< HEAD
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
=======
type UserCreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=30"`
}

type UserUpdateRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=30"`
>>>>>>> d67dbcc (Add all user logic)
}
