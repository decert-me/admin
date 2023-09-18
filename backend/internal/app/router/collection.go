package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitCollectionRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("collection").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("create", api.CreateCollection)                             // 创建合辑
		routersWithAuth.POST("list", api.GetCollectionList)                              // 获取合辑列表
		routersWithAuth.POST("detail", api.GetCollectionDetail)                          // 获取合辑详情
		routersWithAuth.POST("update", api.UpdateCollection)                             // 更新合辑
		routersWithAuth.POST("delete", api.DeleteCollection)                             // 删除合辑
		routersWithAuth.POST("updateStatus", api.UpdateCollectionStatus)                 // 更新合辑状态
		routersWithAuth.POST("updateCollectionQuestSort", api.UpdateCollectionQuestSort) // 编辑合辑下的挑战排序
		routersWithAuth.POST("addQuestToCollection", api.AddQuestToCollection)           // 添加挑战到合辑
		routersWithAuth.POST("collectionQuest", api.GetCollectionQuest)                  // 获取合辑下的挑战
	}
}
