package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"backend/internal/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateTutorial 创建教程
func CreateTutorial(c *gin.Context) {
	var tutorial model.Tutorial
	_ = c.ShouldBindJSON(&tutorial)
	if tutorialBack, err := backend.CreateTutorial(tutorial); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(tutorialBack, "创建成功", c)
	}
}

// GetTutorialList 获取教程列表
func GetTutorialList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindQuery(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := backend.GetTutorialList(pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetTutorial 获取教程详情
func GetTutorial(c *gin.Context) {
	var req request.GetTutorialRequest
	_ = c.ShouldBindQuery(&req)
	if data, err := backend.GetTutorial(req.Id); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}
