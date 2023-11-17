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
		InitUser(db)       // 初始化默认用户
		InitSetting(db)    // 初始化系统设置
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
		model.Collection{},
		model.CollectionRelate{},
		model.Upload{},
		model.SystemSetting{},
	)
	if err != nil {
		global.LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.LOG.Info("register table success")

}

// InitUser 初始化默认账户
func InitUser(db *gorm.DB) {
	// 判断是否存在用户
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		global.LOG.Error("init user failed", zap.Error(err))
		os.Exit(0)
	}
	if count > 0 {
		return
	}
	// 创建角色
	authority := model.Authority{
		AuthorityId:   "888",
		AuthorityName: "超级管理员",
	}
	if err := db.Create(&authority).Error; err != nil {
		global.LOG.Error("create init authority failed", zap.Error(err))
		os.Exit(0)
	}
	// 创建默认用户
	user := model.User{
		Username:    "root",
		AuthorityId: "888",
		Address:     "0xd2AEc55186F9f713128d48087f0e2EF5F453ca79",
	}
	if err := db.Create(&user).Error; err != nil {
		global.LOG.Error("create init user failed", zap.Error(err))
		os.Exit(0)
	}
}

func InitSetting(db *gorm.DB) {
	// 判断是否存在设置
	var count int64
	if err := db.Model(&model.SystemSetting{}).Count(&count).Error; err != nil {
		global.LOG.Error("init system setting failed", zap.Error(err))
		os.Exit(0)
	}
	if count > 0 {
		return
	}
	systemSetting := model.SystemSetting{
		Key:   "beta",
		Value: "true",
	}
	if err := db.Create(&systemSetting).Error; err != nil {
		global.LOG.Error("create init system setting failed", zap.Error(err))
		os.Exit(0)
	}
}
