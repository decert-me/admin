package api

import (
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
)

// SubmitTranslate GitHub Action 提交翻译
func SubmitTranslate(c *gin.Context) {
	var r request.SubmitTranslateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := backend.SubmitTranslate(r); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.Ok(c)
	}
}
