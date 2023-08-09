package model

type AuthorityRelate struct {
	AuthorityId       string `json:"authorityId" gorm:"not null;comment:角色ID"`       // 角色ID
	AuthoritySourceID uint   `json:"authoritySourceId" gorm:"not null;comment:角色ID"` // 资源ID
}
