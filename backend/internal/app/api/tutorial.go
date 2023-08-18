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

// UpdateTutorial 更新教程
func UpdateTutorial(c *gin.Context) {
	var tutorial model.Tutorial
	_ = c.ShouldBindJSON(&tutorial)
	if err := backend.UpdateTutorial(tutorial); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败"+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// UpdateTutorialStatus 更新教程上架状态
func UpdateTutorialStatus(c *gin.Context) {
	var req request.UpdateTutorialStatusRequest
	err := c.ShouldBindJSON(&req)
	if err != nil || req.Status == 0 {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := backend.UpdateTutorialStatus(req.ID, req.Status); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetTutorialList 获取教程列表
func GetTutorialList(c *gin.Context) {
	var pageInfo request.GetTutorialListStatusRequest
	_ = c.ShouldBindJSON(&pageInfo)
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
	_ = c.ShouldBindJSON(&req)
	if data, err := backend.GetTutorial(req.Id); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}

// DeleteTutorial 删除教程
func DeleteTutorial(c *gin.Context) {
	var req request.DelTutorialRequest
	_ = c.ShouldBindJSON(&req)
	if err := backend.DeleteTutorial(req); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
