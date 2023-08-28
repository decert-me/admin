package initialize

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitCommonDB 通用数据库
func InitCommonDB() {
	db := GormPgSql("admin")
	if db != nil {
		global.DB = db
		RegisterTables(db) // 初始化表
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.User{},
		model.UserLogin{},
		model.Authority{},
		model.AuthorityRelate{},
		model.AuthoritySource{},
		model.Tutorial{},
		model.PackLog{},
		model.Category{},
		model.DocType{},
		model.Language{},
		model.Theme{},
		model.Quest{},
		model.Collection{},
	)
	if err != nil {
		global.LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.LOG.Info("register table success")
}
