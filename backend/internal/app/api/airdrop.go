package api

import (
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
)

// RunAirdrop 立即触发空投
func RunAirdrop(c *gin.Context) {
	var req request.RunAirdropReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err = backend.RunAirdrop(req); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.Ok(c)
	}
}

// GetAirdropList 获取空投列表
func GetAirdropList(c *gin.Context) {
	var req request.GetAirdropListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if result, err := backend.GetAirdropList(req); err != nil {
		response.FailWithMessage("获取失败："+err.Error(), c)
	} else {
		response.ResultWithRaw(result, c)
	}
}
