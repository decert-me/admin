package router

import (
	"backend/internal/app/api"
	"github.com/gin-gonic/gin"
)

func InitTranslateRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("translate")
	{
		routersWithAuth.POST("submitTranslate", api.SubmitTranslate) // GitHub Action 提交翻译
	}
}
