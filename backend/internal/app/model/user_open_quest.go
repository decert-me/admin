package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type UserOpenQuest struct {
	gorm.Model
	Address               string         `gorm:"column:address;type:varchar(44);comment:钱包地址;index:address_tokenId" json:"address" form:"address"`
	TokenId               string         `gorm:"column:token_id;index:address_tokenId;type:varchar(100)" json:"token_id"`
	Answer                datatypes.JSON `gorm:"column:answer" json:"answer"`
	OpenQuestScore        int64          `gorm:"column:open_quest_score;default:0;comment:开放题分数" json:"open_quest_score" form:"open_quest_score"`                                       // 开放题分数
	OpenQuestReviewStatus uint8          `gorm:"column:open_quest_review_status;default:0;comment:评阅开放题状态 1 未审核 2 已审核" json:"open_quest_review_status" form:"open_quest_review_status"` //评阅开放题状态 1 未审核 2 已审核
	OpenQuestReviewTime   time.Time      `gorm:"column:open_quest_review_time;comment:评阅开放题时间" json:"open_quest_review_time" form:"open_quest_review_time"`
	Pass                  bool           `gorm:"column:pass;default:false;comment:挑战通过状态" json:"pass" form:"pass"`                          // 状态 false 挑战未通过 true 挑战通过
	UserScore             int64          `gorm:"column:user_score" form:"user_score" json:"user_score"`                                         // 分数
	AiGradingStatus       int8           `gorm:"column:ai_grading_status;default:0;comment:AI判题状态 0-未判题 1-判题中 2-已判题 3-判题失败" json:"ai_grading_status"` // AI判题状态
	AiGradingTime         *time.Time     `gorm:"column:ai_grading_time;comment:AI判题时间" json:"ai_grading_time"`                                  // AI判题时间
}

func (UserOpenQuest) TableName() string {
	return "user_open_quest"
}
