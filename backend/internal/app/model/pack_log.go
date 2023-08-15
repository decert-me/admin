package model

import "gorm.io/gorm"

type PackLog struct {
	gorm.Model
	TutorialID uint
	Status     uint8  `gorm:"column:status;default:1" json:"status,omitempty"` // 状态 1 未打包 2 打包成功 3 打包失败
	Detail     string `gorm:"column:detail;type:text;comment:详情" json:"detail,omitempty"`
}
