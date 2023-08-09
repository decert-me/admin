package model

import (
	"backend/internal/app/global"
)

// UserLogin 用户登录日志
type UserLogin struct {
	global.MODEL
	Username      string `json:"username" form:"username" gorm:"comment:用户登录名"`                                // 用户登录名
	IP            string `json:"ip" form:"ip" gorm:"column:ip;comment:请求ip"`                                   // 请求ip
	LoginLocation string `json:"loginLocation" form:"loginLocation" gorm:"column:login_location;comment:登录地点"` // 登录地点
	OS            string `json:"os" form:"os" gorm:"column:os;comment:操作系统"`                                   // 操作系统
	Browser       string `json:"browser" form:"browser" gorm:"column:browseros;comment:浏览器"`                   // 浏览器
	Status        int    `json:"status" form:"status" gorm:"column:status;comment:登录状态"`                       // 登录状态
	Agent         string `json:"agent" form:"agent" gorm:"column:agent;comment:代理"`                            // 代理
	ErrorMessage  string `json:"error_message" form:"error_message" gorm:"column:error_message;comment:错误信息"`  // 错误信息
}

func (UserLogin) TableName() string {
	return "admin_user_login"
}
