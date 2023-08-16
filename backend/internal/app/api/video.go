package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetYouTubePlayList 获取教程详情
func GetYouTubePlayList(c *gin.Context) {
	var req request.GetYouTubePlayListRequest
	_ = c.ShouldBindJSON(&req)
	if data, err := backend.GetYouTubePlayList(req.Link); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败，"+err.Error(), c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}
