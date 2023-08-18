package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"backend/internal/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetPackList 获取打包列表
func GetPackList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := backend.GetPackList(pageInfo); err != nil {
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

// GetPackLog 获取打包日志
func GetPackLog(c *gin.Context) {
	var req request.GetPackLogRequest
	_ = c.ShouldBindJSON(&req)
	if data, total, err := backend.GetPackLog(req); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     data,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, "获取成功", c)
	}
}

// Pack 打包
func Pack(c *gin.Context) {
	var req request.PackRequest
	_ = c.ShouldBindJSON(&req)
	if err := backend.Pack(req); err != nil {
		global.LOG.Error("打包失败!", zap.Error(err))
		response.FailWithMessage("打包失败"+err.Error(), c)
	} else {
		response.Ok(c)
	}
}
