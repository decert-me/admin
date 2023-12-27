package request

type SubmitTranslateRequest struct {
	Filename string `json:"filename" binding:"required"`
}
