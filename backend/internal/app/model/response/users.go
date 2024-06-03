package response

import "backend/internal/app/model"

type GetUsersListRes struct {
	UserID    uint    `json:"user_id"`
	Address   string  `json:"address"`
	NickName  *string `gorm:"column:nickname;type:varchar(200);default:''" json:"nickname" form:"nickname"`
	Tags      string  `json:"tags"`
	CreatedAt string  `json:"created_at"`
}

type GetUsersInfoRes struct {
	UserID    uint        `json:"user_id"`
	Address   string      `json:"address"`
	NickName  *string     `gorm:"column:nickname;type:varchar(200);default:''" json:"nickname" form:"nickname"`
	CreatedAt string      `json:"created_at"`
	Tag       []model.Tag `gorm:"-" json:"tag"`
}
