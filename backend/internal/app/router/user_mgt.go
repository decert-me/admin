package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserMgtRouter(Router *gin.RouterGroup) {
	routers := Router.Group("userMgt").Use(middleware.JWTAuth())
	{
		routers.POST("register", api.Register) // 注册账号

		routers.POST("resetPassword", api.ResetPassword) // 重置密码
		routers.POST("update", api.UpdateUserInfo)       // 设置用户信息
		// routers.GET("info", api.GetUserInfo)             // 获取用户信息
	}
}
