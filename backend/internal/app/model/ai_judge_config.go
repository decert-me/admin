package model

import (
	"time"
)

// AiJudgeConfig AI判题配置
type AiJudgeConfig struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Config      string    `gorm:"column:config;type:text;not null" json:"config"`
	Enabled     bool      `gorm:"column:enabled;default:false" json:"enabled"`
	AutoGrading bool      `gorm:"column:auto_grading;default:false" json:"auto_grading"` // 是否开启自动判题
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AiJudgeConfig) TableName() string {
	return "ai_judge_config"
}
