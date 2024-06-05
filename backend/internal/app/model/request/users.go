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
	UserID uint    `json:"user_id"`
	Name   *string `gorm:"column:name;type:varchar(200);comment:用户名称;default:''" json:"name" form:"name"`
	TagIds []uint  `gorm:"-" json:"tag_ids"`
}
