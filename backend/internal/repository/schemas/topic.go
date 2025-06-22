package schemas

type CreateTopicRequest struct {
	Title    string `json:"title" binding:"required"`
	ParentID string `json:"parent_id,omitempty"`
}

type TopicResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ParentID string `json:"parent_id"`
}
