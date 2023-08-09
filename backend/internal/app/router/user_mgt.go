package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserMgtRouter(Router *gin.RouterGroup) {
	routers := Router.Group("userMgt").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		routers.GET("list", api.GetUserList)             // 用户列表
		routers.POST("register", api.Register)           // 注册账号
		routers.POST("delete", api.DeleteUser)           // 删除用户
		routers.POST("resetPassword", api.ResetPassword) // 重置密码
		routers.POST("update", api.UpdateUserInfo)       // 设置用户信息
		// routers.GET("info", api.GetUserInfo)             // 获取用户信息
	}
}
