package request

type CreateLabelRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type DeleteLabelRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
