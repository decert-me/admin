package middleware

import (
	"backend/internal/app/global"
	"backend/internal/app/model/response"
	"backend/internal/app/utils"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, err := utils.GetClaims(c)
		if err != nil {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		// 获取请求的PATH
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := waitUse.AuthorityId
		sub2 := waitUse.UUID.String()
		var e *casbin.SyncedEnforcer

		if global.Enforcer != nil {
			e = global.Enforcer
		}

		if sub == "999" { // root admin user
			c.Next()
			return
		}

		// 判断用户是否有被设置过特殊权限
		userPermit := true
		var casbinRules []gormadapter.CasbinRule
		_ = global.DB.Where("ptype = ? AND v0 = ?", "p", sub2).Find(&casbinRules)
		if len(casbinRules) > 0 {
			// 有设置过
			userPermit, _ = e.Enforce(sub2, obj, act)
			if userPermit {
				c.Next()
				return
			} else {
				response.FailWithDetailed(gin.H{}, "权限不足", c)
				c.Abort()
				return
			}
		}

		// 判断角色策略中是否存在
		rolePermit, _ := e.Enforce(sub, obj, act)

		if rolePermit {
			c.Next()
		} else {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
	}
}
