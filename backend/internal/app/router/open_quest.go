package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitOpenQuestRouter(Router *gin.RouterGroup) {
	openQuestPermRouterAuth := Router.Group("/openQuest/perm/").Use(middleware.JWTAuth())
	{
		openQuestPermRouterAuth.POST("getOpenQuestPermList", api.GetOpenQuestPermList) // 获取权限列表
		openQuestPermRouterAuth.POST("addOpenQuestPerm", api.AddOpenQuestPerm)         // 添加权限
		openQuestPermRouterAuth.POST("deleteOpenQuestPerm", api.DeleteOpenQuestPerm)   // 删除权限
	}
}
