package schemas

type UserCreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=30"`
}

type UserUpdateRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=30"`
}
