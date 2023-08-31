package request

import "backend/internal/app/model"

type GetQuestListRequest struct {
	model.Quest
	PageInfo
	Address   string
	OrderKey  string `json:"order_key" form:"order_key,default=token_id"` // 排序
	Desc      bool   `json:"desc" form:"desc,default=true"`               // 排序方式:升序false(默认)|降序true
	SearchKey string `json:"search_key" form:"search_key"`                // 搜索关键字
}

type AddQuestRequest struct {
	Uri         string `json:"uri" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	StartTs     string `json:"start_ts" binding:"required"`
	EndTs       string `json:"end_ts" binding:"required"`
	Supply      string `json:"supply" binding:"required"`
}

type UpdateQuestRequest struct {
	ID           uint  `json:"id"`
	Difficulty   *uint `json:"difficulty"`    // 0:easy;1:moderate;2:difficult
	EstimateTime *uint `json:"estimate_time"` // 预估时间/min
	CollectionID *uint `json:"collection_id"`
	Sort         *int  `json:"sort"` // 排序
}

type UpdateRecommendRequest struct {
	TokenId   int64  `json:"token_id" binding:"required"`
	Recommend string `json:"recommend"` // 推荐
}

type TopQuestRequest struct {
	ID  []uint  `json:"id"`
	Top []*bool `json:"top"`
}

type UpdateQuestStatusRequest struct {
	ID     uint  `json:"id"`
	Status uint8 `json:"status" binding:"required"`
}

type DeleteQuestRequest struct {
	ID uint `json:"id"`
}
