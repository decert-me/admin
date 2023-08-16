package request

type CreateLabelRequest struct {
	Type    string `json:"type" binding:"required"`
	Chinese string `json:"chinese" binding:"required"`
	English string `json:"english" binding:"required"`
	Weight  int    `json:"weight"`
}

type DeleteLabelRequest struct {
	Type string `json:"type" binding:"required"`
	ID   uint   `json:"id" binding:"required"`
}

type GetLabelRequest struct {
	Type string `json:"type" binding:"required"`
}

type UpdateLabelRequest struct {
	Type    string `json:"type" binding:"required"`
	ID      uint   `json:"id" binding:"required"`
	Chinese string `json:"chinese" binding:"required"`
	English string `json:"english" binding:"required"`
	Weight  int    `json:"weight"`
}
