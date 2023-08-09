package model

type AuthoritySource struct {
	ID        uint   `gorm:"primarykey"`                                       // 主键ID
	ModelName string `json:"model_name" form:"model_name" gorm:"comment:模块名称"` // 模块名称
	ModelUrl  string `json:"model_url" form:"model_url" gorm:"comment:模块路由"`   // 模块路由
}
