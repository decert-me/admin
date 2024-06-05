package response

import "backend/internal/app/model"

type GetUsersListRes struct {
	UserID    uint    `json:"user_id"`
	Address   string  `json:"address"`
	Name      *string `gorm:"column:name;type:varchar(200);comment:用户名称;default:''" json:"name" form:"name"`
	Tags      string  `json:"tags"`
	CreatedAt string  `json:"created_at"`
}

type GetUsersInfoRes struct {
	UserID    uint        `json:"user_id"`
	Address   string      `json:"address"`
	Name      *string     `gorm:"column:name;type:varchar(200);comment:用户名称;default:''" json:"name" form:"name"`
	CreatedAt string      `json:"created_at"`
	Tag       []model.Tag `gorm:"-" json:"tag"`
}
