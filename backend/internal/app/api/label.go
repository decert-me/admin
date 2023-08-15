package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateLabel 创建标签
func CreateLabel(c *gin.Context) {
	var label request.CreateLabelRequest
	_ = c.ShouldBindJSON(&label)
	var err error
	if label.Type == "language" {
		err = backend.LabelAddLang(label.Content)
	} else if label.Type == "category" {
		err = backend.LabelAddCategory(label.Content)
	} else if label.Type == "theme" {
		err = backend.LabelAddTheme(label.Content)
	} else {
		response.FailWithMessage("参数错误", c)
	}
	if err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteLabel 删除标签
func DeleteLabel(c *gin.Context) {
	var label request.DeleteLabelRequest
	_ = c.ShouldBindJSON(&label)
	var err error
	if label.Type == "language" {
		err = backend.LabelRemoveLang(label.Content)
	} else if label.Type == "category" {
		err = backend.LabelRemoveCategory(label.Content)
	} else if label.Type == "theme" {
		err = backend.LabelRemoveTheme(label.Content)
	} else {
		response.FailWithMessage("参数错误", c)
	}
	if err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}
