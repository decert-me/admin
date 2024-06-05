package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetTagList 获取标签列表
func GetTagList(c *gin.Context) {
	var r request.GetTagListReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if list, total, err := backend.GetTagList(r); err != nil {
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

// GetTagInfo 获取标签详情
func GetTagInfo(c *gin.Context) {
	var r request.GetTagInfoReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if label, err := backend.GetTagInfo(r); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(label, "获取成功", c)
	}
}

// TagAdd 添加标签
func TagAdd(c *gin.Context) {
	var r model.Tag
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := backend.TagAdd(r); err != nil {
		global.LOG.Error("添加失败!", zap.Error(err))
		response.FailWithErrorMessage("添加失败", err, c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// TagUpdate 修改标签
func TagUpdate(c *gin.Context) {
	var r model.Tag
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := backend.TagUpdate(r); err != nil {
		response.FailWithErrorMessage("更新失败", err, c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetTagUserList 查询标签用户列表
func GetTagUserList(c *gin.Context) {
	var r request.GetTagUserListReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if list, total, err := backend.GetTagUserList(r); err != nil {
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

// TagUserUpdate 添加用户标签
func TagUserUpdate(c *gin.Context) {
	var r request.TagUserUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if err := backend.TagUserUpdate(r); err != nil {
		response.FailWithErrorMessage("更新失败", err, c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// TagDeleteBatch 批量删除标签
func TagDeleteBatch(c *gin.Context) {
	var r request.TagDeleteBatchReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if err := backend.TagDeleteBatch(r); err != nil {
		response.FailWithErrorMessage("删除失败", err, c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// TagUserDeleteBatch 批量删除用户标签
func TagUserDeleteBatch(c *gin.Context) {
	var r request.TagUserDeleteBatchReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if err := backend.TagUserDeleteBatch(r); err != nil {
		response.FailWithErrorMessage("删除失败", err, c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
