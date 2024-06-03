package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetChallengeStatistics 挑战详情统计
func GetChallengeStatistics(c *gin.Context) {
	var r request.GetChallengeStatisticsReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if list, total, err := backend.GetChallengeStatistics(r); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     r.Page,
			PageSize: r.PageSize,
		}, "获取成功", c)
	}
}

// GetChallengeUserStatistics 挑战者统计
func GetChallengeUserStatistics(c *gin.Context) {
	var r request.GetChallengeUserStatisticsReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if list, total, err := backend.GetChallengeUserStatistics(r); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     r.Page,
			PageSize: r.PageSize,
		}, "获取成功", c)
	}
}
