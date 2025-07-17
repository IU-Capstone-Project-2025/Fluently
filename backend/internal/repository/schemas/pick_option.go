package schemas

// CreatePickOptionRequest is a request body for creating a pick option
type CreatePickOptionRequest struct {
	WordID     string   `json:"word_id" validate:"required,uuid"`
	SentenceID string   `json:"sentence_id" validate:"required,uuid"`
	Options    []string `json:"options" validate:"required,len=3,dive,required"`
}

// PickOptionResponse is a response for pick options
type PickOptionResponse struct {
	ID         string   `json:"id"`
	WordID     string   `json:"word_id"`
	SentenceID string   `json:"sentence_id"`
	Options    []string `json:"options"`
}
