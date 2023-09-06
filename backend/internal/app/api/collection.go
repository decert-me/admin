package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateCollection 创建合辑
func CreateCollection(c *gin.Context) {
	var r request.CreateCollectionRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err = backend.CreateCollection(r); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// GetCollectionList 获取合辑列表
func GetCollectionList(c *gin.Context) {
	var r request.GetCollectionListRequest
	_ = c.ShouldBindJSON(&r)
	if list, total, err := backend.GetCollectionList(r); err != nil {
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

// GetCollectionDetail 获取合辑详情
func GetCollectionDetail(c *gin.Context) {
	var r request.GetCollectionDetailRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if detail, err := backend.GetCollectionDetail(r); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithDetailed(detail, "获取成功", c)
	}
}

// UpdateCollection 更新合辑
func UpdateCollection(c *gin.Context) {
	var r request.UpdateCollectionRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err = backend.UpdateCollection(r); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DeleteCollection 删除合辑
func DeleteCollection(c *gin.Context) {
	var r request.DeleteCollectionRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err = backend.DeleteCollection(r); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateCollectionStatus 更新合辑上架状态
func UpdateCollectionStatus(c *gin.Context) {
	var r request.UpdateCollectionStatusRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err = backend.UpdateCollectionStatus(r); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetCollectionQuest 获取合辑内挑战
func GetCollectionQuest(c *gin.Context) {
	var r request.GetCollectionQuestRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if list, err := backend.GetCollectionQuest(r); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}

// UpdateCollectionQuestSort 编辑合辑内挑战排序
func UpdateCollectionQuestSort(c *gin.Context) {
	var r request.UpdateCollectionQuestSortRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err = backend.UpdateCollectionQuestSort(r); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// AddQuestToCollection 添加挑战到合辑内
func AddQuestToCollection(c *gin.Context) {
	var r request.AddQuestToCollectionRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err = backend.AddQuestToCollection(r); err != nil {
		global.LOG.Error("添加失败!", zap.Error(err))
		response.FailWithMessage("添加失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}
