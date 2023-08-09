package global

import (
	"backend/internal/app/config"
	"github.com/allegro/bigcache/v3"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var (
	StartTime  time.Time          // 记录运行时间
	DB         *gorm.DB           // 数据库链接
	CONFIG     config.Server      // 配置信息
	LOG        *zap.Logger        // 日志框架
	TokenCache *bigcache.BigCache // Token 缓存
	Enforcer   *casbin.SyncedEnforcer
)
