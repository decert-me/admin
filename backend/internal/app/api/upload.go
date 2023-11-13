package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"backend/internal/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"path"
	"strings"
)

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	//var file request.FileUploadRequest
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	// 读取文件后缀
	ext := path.Ext(header.Filename)
	ext = strings.ToLower(ext)
	// 限制文件后缀
	if (ext == ".jpg" || ext == ".png" || ext == ".jpeg" || ext == ".gif") == false {
		global.LOG.Error("文件格式不正确", zap.String("ext", ext))
		response.FailWithMessage("文件格式不正确", c)
		return
	}
	// 文件大小限制
	if header.Size > 1024*1024*5 {
		response.FailWithMessage("文件大小超过限制！", c)
		return
	}
	userID := utils.GetUserID(c)
	err, file := backend.UploadAvatar(userID, header) // 文件上传后拿到文件路径
	if err != nil {
		global.LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("上传失败", c)
		return
	}
	response.OkWithDetailed(file, "上传成功", c)
}
