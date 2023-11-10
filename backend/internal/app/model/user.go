package model

import (
	"backend/internal/app/global"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	global.MODEL
	UUID               uuid.UUID `json:"uuid" gorm:"comment:用户UUID"` // 用户UUID
	Username           string    `json:"username" gorm:"comment:用户登录名;unique"`
	Address            string    `json:"address" gorm:"comment:钱包地址;unique"`
	Password           string    `json:"-"  gorm:"comment:用户登录密码"`                  // 用户登录密码
	Nickname           string    `json:"nickname" gorm:"default:系统用户;comment:用户昵称"` // 用户昵称
	HeaderImg          string    `json:"headerImg" gorm:"comment:用户头像"`             // 用户头像
	AuthorityId        string    `json:"-" gorm:"comment:用户角色ID"`                   // 用户角色ID
	Phone              string    `json:"phone"  gorm:"comment:用户手机号"`               // 用户手机号
	Email              string    `json:"email"  gorm:"comment:用户邮箱"`                // 用户邮箱
	AuthoritySourceIds []uint    `json:"authoritySourceIds" gorm:"-"`               // 用户资源ID
	Authority          Authority `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
}
