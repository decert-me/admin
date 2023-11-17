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
