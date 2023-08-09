package request

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int `json:"page" form:"page,default=1"`          // 页码
	PageSize int `json:"pageSize" form:"pageSize,default=10"` // 每页大小
}

type IDRequest struct {
	ID uint `json:"id" form:"id"`
}

type HashRequest struct {
	Hash string `json:"hash" form:"hash"`
}

type PublishItem struct {
	ID        uint  `json:"id" form:"id"`               // 主键ID
	IsPublish *bool `json:"isPublish" form:"isPublish"` // 发布状态
}
