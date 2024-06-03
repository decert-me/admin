package request

type GetTagListReq struct {
	PageInfo
	SearchVal string `form:"search_val" json:"search_val"`
}

type GetTagInfoReq struct {
	TagID uint `form:"tag_id" json:"tag_id"  binding:"required"`
}
type GetTagUserListReq struct {
	PageInfo
	TagID     uint   `form:"tag_id" json:"tag_id"  binding:"required"`
	SearchVal string `form:"search_val" json:"search_val"`
}

type TagUserUpdateReq struct {
	TagID    uint    `form:"tag_id" json:"tag_id"  binding:"required"`
	UserID   uint    `json:"user_id"  binding:"required"`
	NickName *string `json:"nickname" form:"nickname"`
}

type TagDeleteBatchReq struct {
	TagIDs []int64 `json:"tag_ids"  binding:"required"`
}

type TagUserDeleteBatchReq struct {
	TagID   int64   `form:"tag_id" json:"tag_id"  binding:"required"`
	UserIDs []int64 `json:"user_ids"  binding:"required"`
}
