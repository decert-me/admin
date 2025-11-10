package router

import (
	"backend/internal/app/api"
	"github.com/gin-gonic/gin"
)

func InitStatisticsRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("statistics")
	{
		routersWithAuth.POST("getChallengeStatistics", api.GetChallengeStatistics)               // 挑战详情统计
		routersWithAuth.POST("getChallengeUserStatistics", api.GetChallengeUserStatistics)       // 挑战者统计
		routersWithAuth.POST("getChallengeStatisticsSummary", api.GetChallengeStatisticsSummary) // 挑战详情总计
		routersWithAuth.POST("getBootcampChallengeStatistics", api.GetBootcampChallengeStatistics) // 训练营挑战统计
		routersWithAuth.POST("getBootcampChallengeConfig", api.GetBootcampChallengeConfig)       // 获取训练营挑战配置
		routersWithAuth.POST("updateBootcampChallengeConfig", api.UpdateBootcampChallengeConfig) // 更新训练营挑战配置
		routersWithAuth.POST("getEnabledBootcampChallenges", api.GetEnabledBootcampChallenges)   // 获取启用的训练营挑战
	}
}
