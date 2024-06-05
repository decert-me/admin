package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetUsersList 获取用户列表
func GetUsersList(c *gin.Context) {
	var r request.GetUsersListReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if list, total, err := backend.GetUsersList(r); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithErrorMessage("获取失败", err, c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     r.Page,
			PageSize: r.PageSize,
		}, "获取成功", c)
	}
}

// GetUsersInfo 查询用户详情
func GetUsersInfo(c *gin.Context) {
	var r request.GetUsersInfoReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if label, err := backend.GetUsersInfo(r); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithErrorMessage("获取失败", err, c)
	} else {
		response.OkWithDetailed(label, "获取成功", c)
	}
}

// UpdateUsersInfo 修改用户
func UpdateUsersInfo(c *gin.Context) {
	var r request.UpdateUsersInfoReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if err := backend.UpdateUsersInfo(r); err != nil {
		global.LOG.Error("修改失败!", zap.Error(err))
		response.FailWithErrorMessage("修改失败", err, c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}
