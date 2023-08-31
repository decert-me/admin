package api

import (
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
)

// GetQuestList 获取挑战列表
func GetQuestList(c *gin.Context) {
	var r request.GetQuestListRequest
	_ = c.ShouldBindJSON(&r)
	r.Address = c.GetString("address")
	if list, total, err := backend.GetQuestList(r); err != nil {
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

// GetQuest 获取挑战详情
func GetQuest(c *gin.Context) {
	if list, err := backend.GetQuest(c.Param("id")); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithData(list, c)
	}
}

// TopQuest 置顶挑战
func TopQuest(c *gin.Context) {
	var r request.TopQuestRequest
	_ = c.ShouldBindJSON(&r)
	if err := backend.TopQuest(r); err != nil {
		response.FailWithMessage("置顶失败："+err.Error(), c)
	} else {
		response.OkWithMessage("置顶成功", c)
	}
}

// UpdateQuestStatus 修改上架状态
func UpdateQuestStatus(c *gin.Context) {
	var r request.UpdateQuestStatusRequest
	_ = c.ShouldBindJSON(&r)
	if err := backend.UpdateQuestStatus(r); err != nil {
		response.FailWithMessage("更新失败："+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// UpdateQuest 修改挑战
func UpdateQuest(c *gin.Context) {
	var r request.UpdateQuestRequest
	_ = c.ShouldBindJSON(&r)
	if err := backend.UpdateQuest(r); err != nil {
		response.FailWithMessage("更新失败："+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DeleteQuest 删除挑战
func DeleteQuest(c *gin.Context) {
	var r request.DeleteQuestRequest
	_ = c.ShouldBindJSON(&r)
	if err := backend.DeleteQuest(r); err != nil {
		response.FailWithMessage("删除失败："+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
