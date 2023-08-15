package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitPackRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("pack").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("getPackList", api.GetPackList) // 获取打包列表
		routersWithAuth.POST("getPackLog", api.GetPackLog)   // 获取打包日志
		routersWithAuth.POST("pack", api.Pack)               // 打包
	}
}
