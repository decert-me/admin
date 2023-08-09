package model

import "backend/internal/app/global"

type Authority struct {
	global.MODEL
	AuthorityId   string `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID"` // 角色ID
	AuthorityName string `json:"authorityName" gorm:"comment:角色名"`                            // 角色名
	Users         []User `json:"-" gorm:"many2many:user_authority;"`
}
