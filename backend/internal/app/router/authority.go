package router

import (
	"backend/internal/app/api"
	"backend/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

type AuthorityRouter struct{}

func InitAuthorityRouter(Router *gin.RouterGroup) {
	routers := Router.Group("authority").Use(middleware.JWTAuth())
	routersWithAuthority := Router.Group("authority").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		routers.GET("getAuthorityList", api.GetAuthorityList)    // 获取角色列表
		routers.GET("getSourceList", api.GetAuthoritySourceList) // 获取资源列表
		routers.GET("getAuthority", api.GetAuthority)            // 获取角色权限
	}

	{
		routersWithAuthority.POST("createAuthority", api.CreateAuthority)   // 创建角色
		routersWithAuthority.POST("deleteAuthority", api.DeleteAuthority)   // 删除角色
		routersWithAuthority.PUT("updateAuthority", api.UpdateAuthority)    // 更新角色
		routersWithAuthority.POST("setDataAuthority", api.SetDataAuthority) // 设置角色资源权限
	}

}
