package schemas

type CreateTopicRequest struct {
	Title string `json:"title" binding:"required"`
}

type TopicResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
