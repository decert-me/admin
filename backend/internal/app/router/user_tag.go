package router

import (
	"backend/internal/app/api"
	"github.com/gin-gonic/gin"
)

func InitTagRouter(Router *gin.RouterGroup) {
	routersWithAuth := Router.Group("tag")
	{
		routersWithAuth.POST("getTagList", api.GetTagList)                 // 获取标签列表
		routersWithAuth.POST("getTagInfo", api.GetTagInfo)                 // 获取标签详情
		routersWithAuth.POST("tagAdd", api.TagAdd)                         // 添加标签
		routersWithAuth.POST("tagUpdate", api.TagUpdate)                   // 修改标签
		routersWithAuth.POST("getTagUserList", api.GetTagUserList)         // 获取标签用户列表
		routersWithAuth.POST("tagUserUpdate", api.TagUserUpdate)           // 添加用户标签
		routersWithAuth.POST("tagDeleteBatch", api.TagDeleteBatch)         // 批量删除标签
		routersWithAuth.POST("tagUserDeleteBatch", api.TagUserDeleteBatch) // 批量删除用户标签
	}
}
