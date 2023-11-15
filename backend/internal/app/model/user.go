package model

import (
	"backend/internal/app/global"
)

type User struct {
	global.MODEL
	Username    string    `json:"username" gorm:"comment:用户登录名"`
	Address     string    `json:"address" gorm:"comment:钱包地址;unique"`
	HeaderImg   string    `json:"headerImg" gorm:"comment:用户头像"` // 用户头像
	AuthorityId string    `json:"-" gorm:"comment:用户角色ID"`       // 用户角色ID
	Authority   Authority `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
}
