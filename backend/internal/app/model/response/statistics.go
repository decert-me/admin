package response

import "time"

type GetChallengeStatisticsRes struct {
	Address      string    `gorm:"column:address;type:varchar(44);UNIQUE;comment:钱包地址" json:"address" form:"address"`
	Pass         bool      `gorm:"pass" json:"pass"`                   // 是否通过
	FinishTime   time.Time `gorm:"finish_time" json:"finish_time"`     //首次提交时间
	HighestScore int64     `gorm:"highest_score" json:"highest_score"` // 最高分
	Claimed      bool      `gorm:"claimed" json:"claimed"`             // 是否已经领取
	HasDid       bool      `gorm:"has_did" json:"has_did"`             // 已领取链下证书
}
