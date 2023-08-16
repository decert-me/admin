package initialize

import "backend/internal/app/global"

func CheckConfig() {
	if global.CONFIG.Pack.Path == "" || global.CONFIG.Pack.PublishPath == "" {
		panic("pack config error")
	}
}
