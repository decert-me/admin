package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitTutorialRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("tutorial").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("getTutorialList", api.GetTutorialList)           // 获取教程列表
		routersWithAuth.POST("createTutorial", api.CreateTutorial)             // 创建教程
		routersWithAuth.POST("getTutorial", api.GetTutorial)                   // 获取教程详情
		routersWithAuth.POST("deleteTutorial", api.DeleteTutorial)             // 删除教程
		routersWithAuth.POST("updateTutorial", api.UpdateTutorial)             // 更新教程
		routersWithAuth.POST("updateTutorialStatus", api.UpdateTutorialStatus) // 更新教程上架状态

	}
}
