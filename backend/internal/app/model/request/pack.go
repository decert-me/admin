package request

type GetPackLogRequest struct {
	PageInfo
	ID uint `json:"id"`
}

type PackRequest struct {
	ID uint `json:"id"`
}
