package main

import (
	"backend/internal/app/core"
	"backend/internal/app/global"
	"backend/internal/app/initialize"
	"backend/internal/app/service/backend"
	"backend/internal/app/service/system"
	"go.uber.org/zap"
	"time"
)

func main() {
	global.StartTime = time.Now()
	// 初始化Viper
	core.Viper()
	// 初始化zap日志库
	global.LOG = core.Zap()
	// 注册全局logger
	zap.ReplaceGlobals(global.LOG)
	// 初始化数据库
	initialize.InitCommonDB()
	//// 检查配置
	//initialize.CheckConfig()
	// 初始化缓存
	global.TokenCache = initialize.TokenCache()
	// 初始化打包状态
	initialize.InitTutorialPackStatus()
	// 初始化 backend
	backend.Init()
	// 初始化casbin
	system.CasbinInit()
	core.RunWindowsServer()
}
