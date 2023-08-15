package api

import (
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("接收文件失败", c)
		return
	}
	// 文件大小限制
	if header.Size > 1024*1024*20 {
		response.FailWithMessage("文件大小超出限制", c)
		return
	}
	err, hash := backend.IPFSUploadFile(header) // 文件上传后拿到文件路径
	if err != nil {
		response.FailWithMessage("上传失败", c)
		return
	}
	response.OkWithDetailed(response.UploadResponse{Hash: hash}, "上传成功", c)
}
