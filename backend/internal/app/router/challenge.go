package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitChallengeRouter(Router *gin.RouterGroup) {
	challengeRouterAuth := Router.Group("challenge").Use(middleware.JWTAuth())
	{
		challengeRouterAuth.POST("getUserOpenQuestList", api.GetUserOpenQuestList) // 获取用户开放题列表
		challengeRouterAuth.POST("getUserOpenQuest", api.GetUserOpenQuest)         // 获取用户开放题
		challengeRouterAuth.POST("reviewOpenQuest", api.ReviewOpenQuest)           // 审核开放题目
	}
}
