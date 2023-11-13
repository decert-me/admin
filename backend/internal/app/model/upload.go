package model

import "gorm.io/gorm"

// Upload 客户上传文件日志
type Upload struct {
	gorm.Model
	UserID uint   `json:"user_id"`                           // 客户ID
	Name   string `json:"name" gorm:"size:512;comment:文件名"`  // 文件名
	Url    string `json:"url" gorm:"size:2047;comment:文件地址"` // 文件地址
	Tag    string `json:"tag" gorm:"size:30;comment:文件标签"`   // 文件标签
	Key    string `json:"key" gorm:"size:256;comment:编号"`    // 编号
	Type   string `json:"type" gorm:"size:30;comment:文件类型"`  // 文件类型
}
