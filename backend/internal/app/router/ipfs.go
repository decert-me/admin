package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitIPFSRouter(Router *gin.RouterGroup) {
	ipfsRouterAuth := Router.Group("ipfs").Use(middleware.JWTAuth())
	{
		ipfsRouterAuth.POST("uploadFile", api.UploadFile)
	}
}
