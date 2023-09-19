package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAddressInfo 获取地址信息
func GetAddressInfo(c *gin.Context) {
	var req request.GetAddressInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if addressInfo, err := backend.GetAddressInfo(req); err != nil {
		global.LOG.Error("获取地址信息失败!", zap.Error(err))
		response.FailWithMessage("获取地址信息失败："+err.Error(), c)
	} else {
		response.OkWithData(addressInfo, c)
	}
}
