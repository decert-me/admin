package request

import "backend/internal/app/model"

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
type GetTutorialListStatusRequest struct {
	PageInfo
	model.Tutorial
}

type TopTutorialRequest struct {
	ID  []uint  `json:"id"`
	Top []*bool `json:"top"`
}

type UpdateTutorialSortRequest struct {
	ID           uint `json:"id"`
	TutorialSort *int `json:"tutorial_sort"`
}
