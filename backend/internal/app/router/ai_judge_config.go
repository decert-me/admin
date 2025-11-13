package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitAiJudgeConfigRouter(Router *gin.RouterGroup) {
	aiJudgeConfigRouter := Router.Group("ai-judge-config").Use(middleware.JWTAuth())
	{
		aiJudgeConfigRouter.GET("list", api.GetAiJudgeConfigList)                  // 获取配置列表
		aiJudgeConfigRouter.POST("create", api.CreateAiJudgeConfig)                // 创建配置
		aiJudgeConfigRouter.POST("update", api.UpdateAiJudgeConfig)                // 更新配置
		aiJudgeConfigRouter.POST("delete", api.DeleteAiJudgeConfig)                // 删除配置
		aiJudgeConfigRouter.POST("toggle", api.ToggleAiJudgeConfig)                // 切换启用状态
		aiJudgeConfigRouter.GET("enabled", api.GetEnabledAiJudgeConfig)            // 获取当前启用的配置
		aiJudgeConfigRouter.POST("grade", api.AiGrade)                             // AI判题
		aiJudgeConfigRouter.POST("toggle-auto-grading", api.ToggleAutoGrading)    // 切换自动判题状态
		aiJudgeConfigRouter.GET("pending-list", api.GetPendingGradeList)           // 获取待判题列表
		aiJudgeConfigRouter.POST("batch-grade", api.BatchGrade)                    // 批量AI判题（自动提交）
		aiJudgeConfigRouter.GET("history", api.GetAiGradeHistory)                  // 获取AI判题历史
		aiJudgeConfigRouter.POST("batch-grade-preview", api.BatchGradePreview)    // 批量AI判题预览（不提交）
		aiJudgeConfigRouter.POST("submit-batch-grade", api.SubmitBatchGrade)      // 提交批量判题结果
	}
}
