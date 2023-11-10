package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func InitUploadRouter(Router *gin.RouterGroup) {
	routers := Router.Group("upload").Use(middleware.JWTAuth())
	{
		routers.POST("avatar", api.UploadAvatar) // 上传头像
	}
}
