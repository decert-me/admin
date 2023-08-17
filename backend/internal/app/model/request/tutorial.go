package request

type GetTutorialRequest struct {
	Id uint `json:"id"`
}

type DelTutorialRequest struct {
	Id uint `json:"id"`
}

type UpdateTutorialStatusRequest struct {
	ID     uint  `json:"id"`
	Status uint8 `json:"status"` // 状态 1 未上架 2 已上架
}
