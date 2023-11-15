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
		routers.GET("captcha", api.Captcha)                 // 获取验证码
		routers.GET("getLoginMessage", api.GetLoginMessage) // 获取登录签名消息
		routers.POST("authLoginSign", api.AuthLoginSign)    // 校验登录签名
	}
	{
		routersWithAuth.GET("list", api.GetUserList)               // 用户列表
		routersWithAuth.GET("info", api.GetSelfInfo)               // 用户个人资料
		routersWithAuth.POST("update", api.UpdateSelfInfo)         // 用户更改资料
		routersWithAuth.POST("register", api.Register)             // 添加用户
		routersWithAuth.POST("updateUserInfo", api.UpdateUserInfo) // 管理员更改用户资料
		routersWithAuth.POST("delete", api.DeleteUser)             // 删除用户
	}
}
