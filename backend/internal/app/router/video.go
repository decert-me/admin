package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitVideoRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("video").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("getYouTubePlayList", api.GetYouTubePlayList) // 获取YouTube视频列表
	}
}
