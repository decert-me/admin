package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitChallengeRouter(Router *gin.RouterGroup) {
	challengeRouterAuth := Router.Group("challenge").Use(middleware.JWTAuth())
	{
		challengeRouterAuth.POST("getUserOpenQuestList", api.GetUserOpenQuestList)                 // 获取用户开放题列表
		challengeRouterAuth.POST("getUserOpenQuest", api.GetUserOpenQuest)                         // 获取用户开放题
		challengeRouterAuth.POST("reviewOpenQuest", api.ReviewOpenQuest)                           // 审核开放题目
		challengeRouterAuth.POST("getUserOpenQuestListV2", api.GetUserOpenQuestListV2)             // 获取用户开放题列表V2
		challengeRouterAuth.POST("reviewOpenQuestV2", api.ReviewOpenQuestV2)                       // 审核开放题目V2
		challengeRouterAuth.POST("getUserOpenQuestDetailListV2", api.GetUserOpenQuestDetailListV2) // 获取用户开放题详情列表V2
		challengeRouterAuth.POST("getUserQuestDetail", api.GetUserQuestDetail)                     // 获取用户答题详情
	}
}
