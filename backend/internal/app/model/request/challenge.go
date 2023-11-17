package request

import "gorm.io/datatypes"

type GetUserOpenQuestListRequest struct {
	PageInfo
	OpenQuestReviewStatus uint8 `json:"open_quest_review_status" form:"open_quest_review_status"`
}

type GetUserOpenQuestRequest struct {
	ID uint `json:"id" binding:"required"`
}

type ReviewOpenQuestRequest struct {
	ID     uint           `json:"id" binding:"required"`
	Answer datatypes.JSON `json:"answer" binding:"required"`
	Score  int64          `json:"score"`
}