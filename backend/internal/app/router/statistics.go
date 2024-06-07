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
	}
}
