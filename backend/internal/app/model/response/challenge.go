package response

import (
	"backend/internal/app/model"
)

type GetUserOpenQuestListResponse struct {
	Title string `gorm:"column:title;comment:标题;type:varchar" json:"title" form:"title"` // 标题
	model.UserOpenQuest
}

type GetUserOpenQuestResponse struct {
	model.UserOpenQuest
}
