package response

import "backend/internal/app/model"

type GetCollectionListRes struct {
	model.Collection
	ClaimNum     int64 `gorm:"-" json:"claim_num"`     // 铸造数量
	ChallengeNum int64 `gorm:"-" json:"challenge_num"` // 挑战人次
	EstimateTime int64 `gorm:"-" json:"estimate_time"` // 预估时间/min
	QuestNum     int64 `gorm:"-" json:"quest_num"`     // 挑战数量
}
