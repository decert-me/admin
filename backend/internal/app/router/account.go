package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitAccountRouter(Router *gin.RouterGroup) {
	ipfsRouterAuth := Router.Group("account").Use(middleware.JWTAuth())
	{
		ipfsRouterAuth.POST("getAddressInfo", api.GetAddressInfo) // 获取地址信息
	}
}
