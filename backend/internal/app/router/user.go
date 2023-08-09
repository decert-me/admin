package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	routers := Router.Group("user")
	routersWithAuth := Router.Group("user").Use(middleware.JWTAuth())

	{
		routers.POST("login", api.Login)                   // 用户登陆
		routers.GET("captcha", api.Captcha)                // 获取验证码
		routers.POST("changePassword", api.ChangePassword) // 用户修改密码
	}

	{
		routers.GET("info", api.GetSelfInfo)               // 用户个人资料
		routersWithAuth.POST("update", api.UpdateSelfInfo) // 用户更改资料
	}
}
