package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitAirdropRouter(Router *gin.RouterGroup) {
	ipfsRouterAuth := Router.Group("airdrop").Use(middleware.JWTAuth())
	{
		ipfsRouterAuth.POST("runAirdrop", api.RunAirdrop)         // 立即触发空投
		ipfsRouterAuth.POST("getAirdropList", api.GetAirdropList) // 获取空投列表
	}
}
