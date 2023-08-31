package initialize

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"go.uber.org/zap"
)

func InitTutorialPackStatus() {
	err := global.DB.Model(&model.Tutorial{}).Where("pack_status = ?", 1).Update("pack_status", 3).Error
	if err != nil {
		global.LOG.Error("init tutorial pack status failed", zap.Error(err))
		return
	}
}
