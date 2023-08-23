package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitQuestRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("quest").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("list", api.GetQuestList)                   // 获取教程列表
		routersWithAuth.POST("topQuest", api.TopQuest)                   //  置顶教程
		routersWithAuth.POST("updateQuestStatus", api.UpdateQuestStatus) // 更新教程上架状态
		routersWithAuth.POST("update", api.UpdateQuest)                  // 更新教程
	}
}
