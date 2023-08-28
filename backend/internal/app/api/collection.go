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
	if backend.CreateCollection(r) != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败 "+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}
