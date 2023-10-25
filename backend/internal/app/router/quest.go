package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitQuestRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("quest").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("list", api.GetQuestList)                                   // 获取教程列表
		routersWithAuth.GET("/:id", api.GetQuest)                                        // 获取教程详情
		routersWithAuth.POST("topQuest", api.TopQuest)                                   //  置顶教程
		routersWithAuth.POST("updateQuestStatus", api.UpdateQuestStatus)                 // 更新教程上架状态
		routersWithAuth.POST("update", api.UpdateQuest)                                  // 更新教程
		routersWithAuth.POST("delete", api.DeleteQuest)                                  // 删除教程
		routersWithAuth.POST("getQuestCollectionAddList", api.GetQuestCollectionAddList) // 获取待添加到合辑教程列表
	}
}
