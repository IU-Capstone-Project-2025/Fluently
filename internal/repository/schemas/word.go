package schemas

type WordCreateRequest struct {
	Word         string  `json:"word" validate:"required,min=1,max=30"`
	CEFR         *string `json:"ceft" validate:"omitempty,oneof=A1 A2 B1 B2 C1 C2"`
	Translation  *string `json:"translation" validate:"omitempty,max=30"`
	PartOfSpeech string  `json:"part_of_speech" validate:"required, max=30"`
	Context      *string `json:"context" validate:"omitempty,max=100"`
}

type WordCreateResponse struct {
	ID string `json:"id"`
	Word string `json:"word"`
	CEFR *string `json:"cefr,omitempty"`
	Translation *string `json:"translation,omitempty"`
	PartOfSpeech string `json:"part_of_speech"`
	Context *string `json:"context,omitempty"`
}

type WordUpdateRequest struct {
    Word         *string `json:"word" validate:"omitempty,min=1,max=30"`
    CEFR         *string `json:"cefr" validate:"omitempty,oneof=A1 A2 B1 B2 C1 C2"`
    Translation  *string `json:"translation" validate:"omitempty,max=30"`
    PartOfSpeech *string `json:"part_of_speech" validate:"omitempty,max=30"`
    Context      *string `json:"context" validate:"omitempty,max=100"`
}
