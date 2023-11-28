package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"github.com/spf13/cast"
	"time"
)

// GetUserOpenQuestList 获取用户开放题列表
func GetUserOpenQuestList(r request.GetUserOpenQuestListRequest) (list []response.GetUserOpenQuestListResponse, total int64, err error) {
	db := global.DB.Model(&model.UserOpenQuest{})
	db.Select("user_open_quest.*, quest.title").Joins("left join quest on quest.token_id = user_open_quest.token_id")
	if r.OpenQuestReviewStatus != 0 {
		db.Where("open_quest_review_status = ?", r.OpenQuestReviewStatus)
	}
	db.Where("quest.status = 1")
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id desc").Scopes(Paginate(r.Page, r.PageSize)).Find(&list).Error
	return
}

// GetUserOpenQuest 获取用户开放题详情
func GetUserOpenQuest(r request.GetUserOpenQuestRequest) (res response.GetUserOpenQuestResponse, err error) {
	err = global.DB.Model(&model.UserOpenQuest{}).Where("id = ?", r.ID).First(&res).Error
	return
}

// ReviewOpenQuest 审核开放题目
func ReviewOpenQuest(r request.ReviewOpenQuestRequest) (err error) {
	// 获取UserOpenQuest
	var userOpenQuest model.UserOpenQuest
	if err = global.DB.Model(&model.UserOpenQuest{}).Where("id = ? AND open_quest_review_status = 1", r.ID).First(&userOpenQuest).Error; err != nil {
		return errors.New("获取答案失败")
	}
	// 检查是否有变动
	if r.UpdatedAt != nil && !userOpenQuest.UpdatedAt.Equal(*r.UpdatedAt) {
		return errors.New("内容有变动，请重新评分")
	}
	// 获取题目
	var quest model.Quest
	if err = global.DB.Model(&model.Quest{}).Where("token_id = ?", userOpenQuest.TokenId).First(&quest).Error; err != nil {
		return errors.New("获取题目失败")
	}
	// 获取分数
	score, pass, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, r.Answer, quest)
	// 写入审核结果
	err = global.DB.Model(&model.UserOpenQuest{}).Where("id = ? AND open_quest_review_status = 1", r.ID).Updates(&model.UserOpenQuest{
		OpenQuestReviewTime:   time.Now(),
		OpenQuestReviewStatus: 2,
		OpenQuestScore:        score,
		Answer:                r.Answer,
		Pass:                  pass,
	}).Error
	// 写入Message
	var message model.UserMessage
	if pass {
		message = model.UserMessage{
			Title:     "恭喜通过挑战",
			TitleEn:   "Congratulations on passing the challenge!",
			Content:   "你在《" + quest.Title + "》的挑战成绩为 " + cast.ToString(score) + " 分，可领取一枚NFT！",
			ContentEn: "Your score for the challenge \"" + quest.Title + "\" is " + cast.ToString(score) + " points, and you can claim an NFT!",
		}
	} else {
		message = model.UserMessage{
			Title:     "挑战未通过",
			TitleEn:   "Challenge failed",
			Content:   "你在《" + quest.Title + "》的挑战成绩为 " + cast.ToString(score) + " 分，请继续加油吧！",
			ContentEn: "Your score for the challenge \"" + quest.Title + "\" is " + cast.ToString(score) + " points, please continue to working hard.",
		}
	}
	message.TokenId = quest.TokenId
	message.Address = userOpenQuest.Address
	err = global.DB.Model(&model.UserMessage{}).Create(&message).Error
	return
}
