package request

type GetUsersListReq struct {
	PageInfo
	SearchTag     string `form:"search_tag" json:"search_tag"`
	SearchAddress string `form:"search_address" json:"search_address"`
}

type GetUsersInfoReq struct {
	UserID uint `json:"user_id"`
}

type UpdateUsersInfoReq struct {
	UserID   uint    `json:"user_id"`
	NickName *string `gorm:"column:nickname;type:varchar(200);default:''" json:"nickname" form:"nickname"`
	TagIds   []uint  `gorm:"-" json:"tag_ids"`
}
