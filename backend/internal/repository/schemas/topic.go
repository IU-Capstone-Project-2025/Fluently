package schemas

// CreateTopicRequest is a request body for creating a topic
type CreateTopicRequest struct {
	Title    string `json:"title" binding:"required"`
	ParentID string `json:"parent_id,omitempty"`
}

// TopicResponse is a response for a topic
type TopicResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ParentID string `json:"parent_id"`
}

// TopicTitleResponse is a response for a topic that only includes the title
type TopicTitleResponse struct {
	Title string `json:"title"`
}
