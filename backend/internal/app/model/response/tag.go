package response

import (
	"backend/internal/app/model"
	"time"
)

type GetTagListRes struct {
	model.Tag
	UserNum int `json:"userNum"`
}

type GetTagUserListRes struct {
	ID        uint      `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time `json:"createdAt"`            // 创建时间
	NickName  string    `gorm:"column:name;type:varchar(200);default:''" json:"nickname" form:"nickname"`
	Address   string    `json:"address" gorm:"comment:钱包地址;unique"`
}
