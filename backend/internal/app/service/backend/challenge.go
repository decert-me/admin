package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gorm.io/datatypes"
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
		return errors.New("该回答已经评分，请勿重复评分")
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

// GetUserOpenQuestListV2 获取用户开放题列表V2
func GetUserOpenQuestListV2(r request.GetUserOpenQuestListRequest) (list []response.UserOpenQuestJsonElements, total int64, err error) {
	offset := (r.Page - 1) * r.PageSize
	limit := r.PageSize
	db := global.DB.Model(&model.UserOpenQuest{})
	if r.OpenQuestReviewStatus != 0 {
		countSQL := `
				SELECT 
					count(1)
				FROM
					user_open_quest
				JOIN
					jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
				JOIN 
					quest ON quest.token_id = user_open_quest.token_id
				WHERE
					user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest'
		`
		if r.OpenQuestReviewStatus == 2 {
			countSQL += " AND (json_element->>'score' IS NOT NULL OR json_element->>'correct' IS NOT NULL)"
		} else {
			countSQL += " AND json_element->>'score' IS NULL AND json_element->>'correct' IS NULL"
		}
		err = db.Raw(countSQL).Scan(&total).Error
		if err != nil {
			return
		}
		dataSQL := `
				SELECT 
					user_open_quest.id,
					user_open_quest.address,
					user_open_quest.token_id,
					CASE 
						WHEN json_element->>'score' IS NOT NULL OR json_element->>'correct' IS NOT NULL THEN 2
						ELSE 1
					END AS open_quest_review_status,
					json_element->>'open_quest_review_time' AS open_quest_review_time,
					user_open_quest.updated_at,
					(idx::int - 1)  AS index,
					json_element->>'type' AS type,
					json_element->>'value' AS value,
					quest.title AS challenge_title,
					(quest.quest_data->'questions')->(idx::int - 1)->>'title' AS title,
					(quest.quest_data->'questions')->(idx::int - 1)->>'score' AS score,
					(quest.quest_data->'questions')->(idx::int - 1)->>'correct' AS correct,
					json_element AS answer
				FROM
					user_open_quest
				JOIN
					jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
				JOIN 
					quest ON quest.token_id = user_open_quest.token_id
				WHERE
					user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest'
		`
		if r.OpenQuestReviewStatus == 2 {
			dataSQL += " AND (json_element->>'score' IS NOT NULL OR json_element->>'correct' IS NOT NULL)"
		} else {
			dataSQL += " AND json_element->>'score' IS NULL AND json_element->>'correct' IS NULL"
		}
		dataSQL += " ORDER BY id asc OFFSET ? LIMIT ?"
		err = db.Raw(dataSQL, offset, limit).Scan(&list).Error
	} else {
		err = db.Raw(`
				SELECT 
					count(1)
				FROM
					user_open_quest
				JOIN
					jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
				JOIN 
					quest ON quest.token_id = user_open_quest.token_id
				WHERE
					user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest'
		`).Scan(&total).Error
		if err != nil {
			return
		}
		err = db.Raw(`
				SELECT 
					user_open_quest.id,
					user_open_quest.address,
					user_open_quest.token_id,
					CASE 
						WHEN json_element->>'score' IS NOT NULL OR json_element->>'correct' IS NOT NULL THEN 2
						ELSE 1
					END AS open_quest_review_status,
					json_element->>'open_quest_review_time' AS open_quest_review_time,
					user_open_quest.updated_at,
					(idx::int - 1) AS index,
					json_element->>'type' AS type,
					json_element->>'value' AS value,
					quest.title AS challenge_title,
					(quest.quest_data->'questions')->(idx::int - 1)->>'title' AS title,
					(quest.quest_data->'questions')->(idx::int - 1)->>'score' AS score,
					(quest.quest_data->'questions')->(idx::int - 1)->>'correct' AS correct,
					json_element AS answer
				FROM
					user_open_quest
				JOIN
					jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
				JOIN 
					quest ON quest.token_id = user_open_quest.token_id
				WHERE
					user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest'
				ORDER BY id asc
				OFFSET ? LIMIT ?
		`, offset, limit).Scan(&list).Error
	}
	return
}

func ReviewOpenQuestV2(req []request.ReviewOpenQuestRequestV2) (err error) {
	// 开启事务
	db := global.DB.Begin()
	// 用户开放题
	userOpenQuestTimeMap := make(map[uint]time.Time)
	// 题目
	questMap := make(map[int64]model.Quest)
	for _, r := range req {
		var userOpenQuest model.UserOpenQuest
		if err = db.Model(&model.UserOpenQuest{}).Where("id = ? AND open_quest_review_status = 1", r.ID).First(&userOpenQuest).Error; err != nil {
			db.Rollback()
			return errors.New("该回答已经评分，请勿重复评分")
		}
		// 如果不存在
		if _, ok := userOpenQuestTimeMap[r.ID]; !ok {
			userOpenQuestTimeMap[r.ID] = userOpenQuest.UpdatedAt
		}
		// 检查是否有变动，跳过
		if r.UpdatedAt != nil && !userOpenQuestTimeMap[r.ID].Equal(*r.UpdatedAt) {
			continue
		}
		// 如果不存在
		if _, ok := questMap[userOpenQuest.TokenId]; !ok {
			var quest model.Quest
			if err = db.Model(&model.Quest{}).Where("token_id = ?", userOpenQuest.TokenId).First(&quest).Error; err != nil {
				db.Rollback()
				return errors.New("获取题目失败")
			}
			questMap[userOpenQuest.TokenId] = quest
		}
		// 获取题目
		quest := questMap[userOpenQuest.TokenId]
		// 填入原有答案
		answer := userOpenQuest.Answer
		answerRes, err := sjson.Set(string(answer), fmt.Sprintf("%d", r.Index), r.Answer)
		if err != nil {
			db.Rollback()
			return errors.New("写入审核结果失败")
		}
		// 判断所有开放题是否审核完成
		var openQuestReviewStatus uint8 = 2 // 已审核
		for _, v := range gjson.Get(answerRes, "@this").Array() {
			// 跳过不是开放题
			if v.Get("type").String() != "open_quest" {
				continue
			}
			// 判断分数是否为空
			if v.Get("score").String() == "" {
				openQuestReviewStatus = 1 // 未审核
				break
			}
		}
		var openQuestReviewTime time.Time // 审核时间
		var pass bool                     // 是否通过
		var userReturnScore int64         // 分数
		if openQuestReviewStatus == 2 {
			openQuestReviewTime = time.Now()
		}
		// 判断是否通过
		if openQuestReviewStatus == 2 {
			userReturnScore, pass, err = AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(answerRes), quest)
			if err != nil {
				db.Rollback()
				return errors.New("服务器错误")
			}
		}
		score := fmt.Sprintf("%d", userReturnScore)
		// 写入审核结果
		err = db.Model(&model.UserOpenQuest{}).Where("id = ? AND open_quest_review_status = 1", r.ID).Updates(&model.UserOpenQuest{
			OpenQuestReviewTime:   openQuestReviewTime,
			OpenQuestReviewStatus: openQuestReviewStatus,
			OpenQuestScore:        userReturnScore,
			Answer:                datatypes.JSON(answerRes),
			Pass:                  pass,
		}).Error
		if err != nil {
			db.Rollback()
			return errors.New("写入结果失败")
		}
		// 审核完成发送消息
		if openQuestReviewStatus == 2 {
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
			err = db.Model(&model.UserMessage{}).Create(&message).Error
			if err != nil {
				db.Rollback()
				return errors.New("发送消息失败")
			}
		}
	}
	return db.Commit().Error
}
