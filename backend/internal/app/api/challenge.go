package api

import (
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
)

// GetUserOpenQuestList 获取用户开放题列表
func GetUserOpenQuestList(c *gin.Context) {
	var r request.GetUserOpenQuestListRequest
	_ = c.ShouldBindJSON(&r)
	if list, total, err := backend.GetUserOpenQuestList(r); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     r.Page,
			PageSize: r.PageSize,
		}, "获取成功", c)
	}

}

// GetUserOpenQuest 获取用户开放题详情
func GetUserOpenQuest(c *gin.Context) {
	var r request.GetUserOpenQuestRequest
	_ = c.ShouldBindJSON(&r)
	if data, err := backend.GetUserOpenQuest(r); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

// ReviewOpenQuest 审核开放题目
func ReviewOpenQuest(c *gin.Context) {
	var r request.ReviewOpenQuestRequest
	_ = c.ShouldBindJSON(&r)
	if err := backend.ReviewOpenQuest(r); err != nil {
		response.FailWithMessage("操作失败："+err.Error(), c)
	} else {
		response.OkWithMessage("操作成功", c)
	}
}

// GetUserOpenQuestListV2 获取用户开放题列表V2
func GetUserOpenQuestListV2(c *gin.Context) {
	var r request.GetUserOpenQuestListRequest
	_ = c.ShouldBindJSON(&r)
	type Detail struct {
		List          interface{} `json:"list"`
		Total         int64       `json:"total"`
		Page          int         `json:"page"`
		PageSize      int         `json:"pageSize"`
		TotalToReview int64       `json:"total_to_review"`
	}
	if list, total, totalToReview, err := backend.GetUserOpenQuestListV2(r); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithDetailed(Detail{
			List:          list,
			Total:         total,
			Page:          r.Page,
			PageSize:      r.PageSize,
			TotalToReview: totalToReview,
		}, "获取成功", c)
	}

}

// ReviewOpenQuestV2 审核开放题目V2
func ReviewOpenQuestV2(c *gin.Context) {
	var r []request.ReviewOpenQuestRequestV2
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if err := backend.ReviewOpenQuestV2(r); err != nil {
		response.FailWithMessage("操作失败："+err.Error(), c)
	} else {
		response.OkWithMessage("操作成功", c)
	}
}

// GetUserOpenQuestDetailListV2 获取用户开放题详情列表
func GetUserOpenQuestDetailListV2(c *gin.Context) {
	var r request.GetUserOpenQuestDetailListRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if list, total, err := backend.GetUserOpenQuestDetailListV2(r); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     r.Page,
			PageSize: r.PageSize,
		}, "获取成功", c)
	}
}

// GetUserQuestDetail 获取用户题目详情
func GetUserQuestDetail(c *gin.Context) {
	var r request.GetUserQuestDetailRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}
	if data, err := backend.GetUserQuestDetail(r); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}
