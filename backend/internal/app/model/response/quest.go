package response

import (
	"backend/internal/app/model"
)

type GetQuestListRes struct {
	model.Quest
	ClaimNum     int64 `json:"claim_num"`     // 铸造数量
	ChallengeNum int64 `json:"challenge_num"` // 挑战人次
}

type ChallengeUsers struct {
	ID          uint    `gorm:"primarykey" json:"-"`
	Address     string  `gorm:"column:address;type:char(42);UNIQUE;comment:钱包地址" json:"address" form:"address"`
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
	Claimed bool `gorm:"claimed" json:"claimed"`
}
