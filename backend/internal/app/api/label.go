package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateLabel 创建标签
func CreateLabel(c *gin.Context) {
	var label request.CreateLabelRequest
	err := c.ShouldBindJSON(&label)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if label.Type == "language" {
		err = backend.LabelAddLang(model.Language{Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else if label.Type == "category" {
		err = backend.LabelAddCategory(model.Category{Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else if label.Type == "theme" {
		err = backend.LabelAddTheme(model.Theme{Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else if label.Type == "challenge" {
		err = backend.LabelAddQuest(model.QuestCategory{Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteLabel 删除标签
func DeleteLabel(c *gin.Context) {
	var label request.DeleteLabelRequest
	err := c.ShouldBindJSON(&label)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if label.Type == "language" {
		err = backend.LabelRemoveLang(model.Language{
			Model: gorm.Model{ID: label.ID},
		})
	} else if label.Type == "category" {
		err = backend.LabelRemoveCategory(model.Category{
			Model: gorm.Model{ID: label.ID},
		})
	} else if label.Type == "theme" {
		err = backend.LabelRemoveTheme(model.Theme{
			Model: gorm.Model{ID: label.ID},
		})
	} else if label.Type == "challenge" {
		err = backend.LabelRemoveQuest(model.QuestCategory{
			Model: gorm.Model{ID: label.ID},
		})
	} else {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败："+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// GetLabelList 获取标签列表
func GetLabelList(c *gin.Context) {
	var label request.GetLabelRequest
	err := c.ShouldBindJSON(&label)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	var data interface{}
	if label.Type == "language" {
		data, err = backend.LabelLangList()
	} else if label.Type == "category" {
		data, err = backend.LabelCategoryList()
	} else if label.Type == "theme" {
		data, err = backend.LabelThemeList()
	} else if label.Type == "challenge" {
		data, err = backend.LabelQuestList()
	} else {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}

// UpdateLabel 修改标签
func UpdateLabel(c *gin.Context) {
	var label request.UpdateLabelRequest
	err := c.ShouldBindJSON(&label)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if label.Type == "language" {
		err = backend.LabelUpdateLang(model.Language{
			Model:   gorm.Model{ID: label.ID},
			Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else if label.Type == "category" {
		err = backend.LabelUpdateCategory(model.Category{
			Model:   gorm.Model{ID: label.ID},
			Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else if label.Type == "theme" {
		err = backend.LabelUpdateTheme(model.Theme{
			Model:   gorm.Model{ID: label.ID},
			Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else if label.Type == "challenge" {
		err = backend.LabelUpdateQuest(model.QuestCategory{
			Model:   gorm.Model{ID: label.ID},
			Chinese: label.Chinese, English: label.English, Weight: label.Weight})
	} else {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}
