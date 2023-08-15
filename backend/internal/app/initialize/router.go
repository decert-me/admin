package initialize

import (
	"backend/internal/app/global"
	"backend/internal/app/middleware"
	"backend/internal/app/router"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 初始化总路由

func Routers() *gin.Engine {
	var Router *gin.Engine
	// 开发环境打开日志 && 打开pprof
	if global.CONFIG.System.Env == "develop" {
		Router = gin.Default()
		pprof.Register(Router) // 性能
	} else {
		Router = gin.New()
		Router.Use(gin.Recovery())
	}

	Router.StaticFile("/favicon.png", "./frontend/build/favicon.png")
	Router.Static("/static", "./frontend/build/static")   // dist里面的静态资源
	Router.StaticFile("/", "./frontend/build/index.html") // 前端网页入口页面
	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面4行注释
	// Router.LoadHTMLGlob("./dist/*.html") // npm打包成dist的路径
	Router.StaticFS(global.CONFIG.Local.Path, http.Dir(global.CONFIG.Local.Path)) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")

	global.LOG.Info("use middleware logger")

	// 跨域，如需跨域可以打开下面的注释
	Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	//Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	global.LOG.Info("use middleware cors")

	// 开发环境打开swagger
	if global.CONFIG.System.Env == "develop" {
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		global.LOG.Info("register swagger handler")
	}
	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		// 状态监测
	}
	v1Group := Router.Group("v1")
	{
		router.InitUserRouter(v1Group)      //用户
		router.InitAuthorityRouter(v1Group) //权限
		router.InitUserMgtRouter(v1Group)
		router.InitTutorialRouter(v1Group)
		router.InitLabelRouter(v1Group)
		router.InitIPFSRouter(v1Group)
		router.InitPackRouter(v1Group)
	}
	global.LOG.Info("router register success")
	return Router
}
