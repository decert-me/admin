package response

import (
	"backend/internal/app/model"
	"time"
)

type GetQuestListRes struct {
	model.Quest
	ClaimNum         int64  `gorm:"-" json:"claim_num"`          // 铸造数量
	ChallengeNum     int64  `gorm:"-" json:"challenge_num"`      // 挑战人次
	ChallengeUserNum int64  `gorm:"-" json:"challenge_user_num"` // 挑战人数
	CollectionID     []uint `gorm:"-" json:"collection_id"`      // 合辑ID
}

type ChallengeUsers struct {
	ID          uint    `gorm:"primarykey" json:"-"`
	Address     string  `gorm:"column:address;type:varchar(44);UNIQUE;comment:钱包地址" json:"address" form:"address"`
	NickName    *string `gorm:"column:nickname;type:varchar(200);default:''" json:"nickname" form:"nickname"`
	Avatar      *string `gorm:"column:avatar;type:varchar(200);comment:用户头像;default:''" json:"avatar" form:"avatar"`
	Description *string `gorm:"column:description;type:varchar(100);comment:自我介绍;default:''" json:"description" form:"description"`
}

type GetQuestChallengeUserRes struct {
	Users []ChallengeUsers `gorm:"users" json:"users"`
	Times int64
}

type GetQuestRes struct {
	model.Quest
	CollectionID []uint `gorm:"-" json:"collection_id"` // 合辑ID
	//Claimed bool `gorm:"claimed" json:"claimed"`
}

type GetQuestCollectionAddListRes struct {
	model.Quest
	CollectionID []uint `gorm:"-" json:"collection_id"` // 合辑ID
}

type GetQuestStatisticsRes struct {
	UUID          string    `gorm:"uuid" json:"uuid"`
	Title         string    `gorm:"title" json:"title"`
	Address       string    `gorm:"column:address;type:varchar(44);UNIQUE;comment:钱包地址" json:"address" form:"address"`
	Name          string    `gorm:"name" json:"name"`
	Tags          string    `json:"tags"`
	ChallengeTime time.Time `gorm:"challenge_time" json:"challenge_time"` // 挑战时间
	QuestID       int64     `gorm:"quest_id" json:"-"`
	TokenID       string    `gorm:"token_id" json:"-"`
	Pass          bool      `gorm:"pass" json:"pass"`
	Claimed       bool      `gorm:"claimed" json:"claimed"`
	ScoreDetail   string    `gorm:"column:score_detail" json:"score_detail"`
	Annotation    string    `gorm:"column:annotation" json:"annotation"` // 批注
}

type GetChallengeUserStatisticsRes struct {
	UserID      int64  `gorm:"user_id" json:"user_id"`
	Address     string `gorm:"column:address;type:varchar(44);UNIQUE;comment:钱包地址" json:"address" form:"address"`
	Name        string `gorm:"name" json:"name"`
	Tags        string `json:"tags"`
	SuccessNum  int64  `json:"success_num"`   // 挑战成功数量
	FailNum     int64  `json:"fail_num"`      // 挑战失败数量
	ClaimNum    int64  `json:"claim_num"`     // 领取NFT数量
	NotClaimNum int64  `json:"not_claim_num"` // 未领取NFT数量
}

type GetQuestStatisticsSummaryRes struct {
	// 挑战数量
	ChallengeNum int64 `json:"challenge_num"`
	// 挑战人数
	ChallengeUserNum int64 `json:"challenge_user_num"`
	// 挑战成功数量
	SuccessNum int64 `json:"success_num"`
	// 挑战失败数量
	FailNum int64 `json:"fail_num"`
	// 领取NFT数量
	ClaimNum int64 `json:"claim_num"`
	// 未领取NFT数量
	NotClaimNum int64 `json:"not_claim_num"`
}
