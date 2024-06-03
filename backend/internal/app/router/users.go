package router

import (
	"backend/internal/app/api"
	"github.com/gin-gonic/gin"
)

func InitUsersRouter(Router *gin.RouterGroup) {
	routers := Router.Group("users")
	{
		routers.POST("getUsersList", api.GetUsersList)       // 获取用户列表
		routers.POST("getUsersInfo", api.GetUsersInfo)       // 查询用户详情
		routers.POST("updateUsersInfo", api.UpdateUsersInfo) // 更新用户信息
	}
}
