package schemas

type DayWordResponse struct {
	Word          string     `json:"word"`
	Translation   string     `json:"translation"`
	Transcription *string     `json:"transcription,omitempty"`
	CEFRLevel     string     `json:"cefr_level"`
	IsLearned     bool       `json:"is_learned"`
	Topic         string     `json:"topic"`
	Subtopic      string     `json:"subtopic"`
	Sentences     []Sentence `json:"sentences"`
	Exercise      Exercise   `json:"exercise"`
}
