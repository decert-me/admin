package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitQuestRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("quest").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("list", api.GetQuestList)                                   // 获取挑战列表
		routersWithAuth.GET("/:id", api.GetQuest)                                        // 获取挑战详情
		routersWithAuth.POST("topQuest", api.TopQuest)                                   //  置顶挑战
		routersWithAuth.POST("updateQuestStatus", api.UpdateQuestStatus)                 // 更新挑战上架状态
		routersWithAuth.POST("update", api.UpdateQuest)                                  // 更新挑战
		routersWithAuth.POST("delete", api.DeleteQuest)                                  // 删除挑战
		routersWithAuth.POST("getQuestCollectionAddList", api.GetQuestCollectionAddList) // 获取待添加到合辑挑战列表
	}
}
