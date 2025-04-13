package schemas

<<<<<<< HEAD
type CreateWordRequest struct {
	Word         string  `json:"word" binding:"required"`
	Translation  *string `json:"translation"`
	PartOfSpeech string  `json:"part_of_speech" binding:"required"`
	Context      *string `json:"context"`
}

type WordResponse struct {
	ID           string  `json:"id"`
	Word         string  `json:"word"`
	Translation  *string `json:"translation,omitempty"`
	PartOfSpeech string  `json:"part_of_speech"`
	Context      *string `json:"context,omitempty"`
=======
type WordCreateRequest struct {
	Word         string  `json:"word" validate:"required,min=1,max=30"`
	CEFR         *string `json:"ceft" validate:"omitempty,oneof=A1 A2 B1 B2 C1 C2"`
	Translation  *string `json:"translation" validate:"omitempty,max=30"`
	PartOfSpeech string  `json:"part_of_speech" validate:"required, max=30"`
	Context      *string `json:"context" validate:"omitempty,max=100"`
>>>>>>> 514fbe1 (Add word create logic)
}
