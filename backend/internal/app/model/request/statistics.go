package request

type GetChallengeStatisticsReq struct {
	PageInfo
	SearchQuest   string `form:"search_quest" json:"search_quest"`
	SearchTag     string `form:"search_tag" json:"search_tag"`
	SearchAddress string `form:"search_address" json:"search_address"`
	Pass          *bool  `gorm:"pass" json:"pass"`
	Claimed       *bool  `gorm:"claimed" json:"claimed"`
}

type GetChallengeUserStatisticsReq struct {
	PageInfo
	SearchTag     string `form:"search_tag" json:"search_tag"`
	SearchAddress string `form:"search_address" json:"search_address"`
}
