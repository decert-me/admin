package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitLabelRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("label").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("delete", api.DeleteLabel) // 删除标签
		routersWithAuth.POST("create", api.CreateLabel) // 创建标签
	}
}
