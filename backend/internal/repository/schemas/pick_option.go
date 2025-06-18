package schemas

type CreatePickOptionRequest struct {
	WordID     string   `json:"word_id" validate:"required,uuid"`
	SentenceID string   `json:"sentence_id" validate:"required,uuid"`
	Options    []string `json:"options" validate:"required,len=3,dive,required"`
}

type PickOptionResponse struct {
	ID         string   `json:"id"`
	WordID     string   `json:"word_id"`
	SentenceID string   `json:"sentence_id"`
	Options    []string `json:"options"`
}
