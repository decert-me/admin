package initialize

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	uuid "github.com/satori/go.uuid"
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
		model.CollectionRelate{},
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
	// 创建默认用户
	user := model.User{
		UUID:        uuid.NewV4(),
		Username:    "root",
		Password:    "$2a$10$7iTyC7BlYSofctYy6.6bq.sL12FfybC/hZQ5K/5lqhuINjAAGSSw2",
		Nickname:    "root",
		AuthorityId: "999",
	}
	if err := db.Create(&user).Error; err != nil {
		global.LOG.Error("create init user failed", zap.Error(err))
		os.Exit(0)
	}
}
