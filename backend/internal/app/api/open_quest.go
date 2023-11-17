package api

import (
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
)

// GetOpenQuestPermList 获取权限列表
func GetOpenQuestPermList(c *gin.Context) {
	var r request.GetOpenQuestPermListRequest
	_ = c.ShouldBindJSON(&r)
	if list, total, err := backend.GetOpenQuestPermList(r); err != nil {
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

// AddOpenQuestPerm 添加权限
func AddOpenQuestPerm(c *gin.Context) {
	var r request.AddOpenQuestPermRequest
	_ = c.ShouldBindJSON(&r)
	if err := backend.AddOpenQuestPerm(r); err != nil {
		response.FailWithMessage("添加失败："+err.Error(), c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// DeleteOpenQuestPerm 删除权限
func DeleteOpenQuestPerm(c *gin.Context) {
	var r request.DeleteOpenQuestPermRequest
	_ = c.ShouldBindJSON(&r)
	if err := backend.DeleteOpenQuestPerm(r); err != nil {
		response.FailWithMessage("删除失败："+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
