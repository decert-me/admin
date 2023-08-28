package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitCollectionRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("collection").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("create", api.CreateCollection) //
	}
}
