package schemas

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
}
