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
