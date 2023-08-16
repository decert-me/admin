package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitLabelRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("label").Use(middleware.JWTAuth())
	{
		routersWithAuth.POST("deleteLabel", api.DeleteLabel)   // 删除标签
		routersWithAuth.POST("createLabel", api.CreateLabel)   // 创建标签
		routersWithAuth.POST("getLabelList", api.GetLabelList) // 获取标签列表
		routersWithAuth.POST("updateLabel", api.UpdateLabel)   // 更新标签
	}
}
