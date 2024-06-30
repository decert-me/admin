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
	"sort"
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
	_, _, score, userScore, pass, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, r.Answer, quest)
	// 写入审核结果
	err = global.DB.Model(&model.UserOpenQuest{}).Where("id = ? AND open_quest_review_status = 1", r.ID).Updates(&model.UserOpenQuest{
		OpenQuestReviewTime:   time.Now(),
		OpenQuestReviewStatus: 2,
		OpenQuestScore:        score,
		Answer:                r.Answer,
		Pass:                  pass,
		UserScore:             userScore,
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
func GetUserOpenQuestListV2(r request.GetUserOpenQuestListRequest) (list []response.UserOpenQuestJsonElements, total int64, totalToReview int64, err error) {
	offset := (r.Page - 1) * r.PageSize
	limit := r.PageSize
	db := global.DB.Model(&model.UserOpenQuest{})
	dataSQL := `
		SELECT
			t.json_element ->> 'title' as title,quest.title as challenge_title,quest.uuid,quest.token_id,(idx::int - 1)  AS index,quest.add_ts as add_ts
		FROM
			quest,
			jsonb_array_elements (quest.quest_data -> 'questions') WITH ORDINALITY AS t(json_element, idx)
		WHERE
			t.json_element ->> 'type' = 'open_quest' AND quest.status = 1
	`
	err = db.Raw(dataSQL).Scan(&list).Error
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		// 待评分数量
		toReviewCountSQL := `
		SELECT 
			count(1)
		FROM
			user_open_quest
		JOIN
			jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
		JOIN 
			quest ON quest.token_id = user_open_quest.token_id
		WHERE
			user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest' AND quest.token_id = ? AND idx= ?
			AND json_element->>'score' IS NULL AND json_element->>'correct' IS NULL  AND json_element->>'score' IS NULL AND json_element->>'correct' IS NULL
		`
		err = global.DB.Raw(toReviewCountSQL, list[i].TokenId, list[i].Index+1).Scan(&list[i].ToReviewCount).Error
		if err != nil {
			continue
		}
		// 已评分数量
		reviewedCountSQL := `
		SELECT 
			count(1)
		FROM
			user_open_quest
		JOIN
			jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
		JOIN 
			quest ON quest.token_id = user_open_quest.token_id
		WHERE
			user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest' AND quest.token_id = ? AND idx= ?
			AND (json_element->>'score' IS NOT NULL OR json_element->>'correct' IS NOT NULL)
		`
		err = global.DB.Raw(reviewedCountSQL, list[i].TokenId, list[i].Index+1).Scan(&list[i].ReviewedCount).Error
		if err != nil {
			continue
		}
		// 查询最新提交时间
		err = global.DB.Model(&model.UserOpenQuest{}).Select("created_at").Where("token_id = ?", list[i].TokenId).Order("id desc").First(&list[i].LastSummitTime).Error
		if err != nil {
			continue
		}
		// 查询最新审核时间
		selectSQL := `
		SELECT
			to_timestamp(t.json_element ->> 'open_quest_review_time','YYYY-MM-DD HH24:MI:SS')
		FROM
			user_open_quest,
			jsonb_array_elements (user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx)
		WHERE
		user_open_quest.token_id = ? AND idx= ? AND t.json_element ->> 'open_quest_review_time' != ''
		ORDER BY t.json_element ->> 'open_quest_review_time' desc
		limit 1
		`
		err = global.DB.Model(&model.UserOpenQuest{}).Raw(selectSQL, list[i].TokenId, list[i].Index+1).Scan(&list[i].LastReviewTime).Error
		if err != nil {
			continue
		}
		fmt.Println("add_ts", list[i].Addts)
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Addts > list[j].Addts
	})
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].ToReviewCount > 0
	})
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].ToReviewCount != 0 {
			if list[i].TokenId == list[j].TokenId {
				return list[i].Index < list[j].Index
			}
			return list[i].Addts > list[j].Addts
		}
		return false
	})
	// 过滤
	temp := make([]response.UserOpenQuestJsonElements, 0)
	for _, v := range list {
		if v.LastSummitTime.IsZero() {
			continue
		}
		temp = append(temp, v)
	}
	for i := 0; i < len(list); i++ {
		// 先按照ToReviewCount倒序排序
		totalToReview += list[i].ToReviewCount
	}
	total = int64(len(temp))
	// limit offset
	result := make([]response.UserOpenQuestJsonElements, 0)

	for i := offset; i < (offset+limit) && i < len(temp); i++ {
		result = append(result, temp[i])
	}
	return result, total, totalToReview, nil
}

// GetUserOpenQuestDetailListV2 获取用户开放题详情
func GetUserOpenQuestDetailListV2(r request.GetUserOpenQuestDetailListRequest) (list []response.GetUserOpenQuestDetailListV2, total int64, err error) {
	offset := (r.Page - 1) * r.PageSize
	limit := r.PageSize
	db := global.DB.Model(&model.UserOpenQuest{})
	// OpenQuestReviewStatus 1 未审核 2 已审核
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
			user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest' AND quest.token_id = ? AND idx= ?
	`
	if r.OpenQuestReviewStatus == 2 {
		countSQL += " AND (json_element->>'score' IS NOT NULL OR json_element->>'correct' IS NOT NULL)"
	} else {
		countSQL += " AND json_element->>'score' IS NULL AND json_element->>'correct' IS NULL"
	}
	err = db.Raw(countSQL, r.TokenID, *r.Index+1).Scan(&total).Error
	if err != nil {
		return
	}
	dataSQL := `
				SELECT 
					user_open_quest.id,
					user_open_quest.address,
					quest.uuid,
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
					json_element AS answer,
					quest.quest_data->>'passingScore' AS pass_score,
					quest.quest_data AS quest_data,
					quest.meta_data AS meta_data,  
					user_open_quest.answer AS user_answer
				FROM
					user_open_quest
				JOIN
					jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
				JOIN 
					quest ON quest.token_id = user_open_quest.token_id
				WHERE
					user_open_quest.deleted_at IS NULL AND quest.status = 1 AND json_element->>'type' = 'open_quest' AND quest.token_id = ? AND idx= ?
		`
	if r.OpenQuestReviewStatus == 2 {
		dataSQL += " AND (json_element->>'score' IS NOT NULL OR json_element->>'correct' IS NOT NULL)"
	} else {
		dataSQL += " AND json_element->>'score' IS NULL AND json_element->>'correct' IS NULL"
	}
	dataSQL += " ORDER BY updated_at asc OFFSET ? LIMIT ?"
	err = db.Raw(dataSQL, r.TokenID, *r.Index+1, offset, limit).Scan(&list).Error
	if err != nil {
		return
	}
	// 计算分数
	for i := 0; i < len(list); i++ {
		quest := model.Quest{
			TokenId:   list[i].TokenId,
			MetaData:  list[i].MetaData,
			QuestData: list[i].QuestData,
		}
		list[i].TotalScore, list[i].UserScore, _, _, _, err = AnswerCheck(global.CONFIG.Quest.EncryptKey, list[i].UserAnswer, quest)
		if err != nil {
			return
		}
		var showStr string
		showStr = fmt.Sprintf("%s...%s", list[i].Address[:6], list[i].Address[len(list[i].Address)-4:])
		// 显示标签
		nickname, name, tags, err := GetUserNameTagsByAddress(list[i].Address)
		if err == nil {
			if nickname != "" {
				showStr = nickname
			}
			if name != "" {
				showStr += "-" + name
			}
			for i := 0; i < len(tags); i++ {
				showStr += "，" + tags[i]
			}
		}

		list[i].NickName = &showStr
	}
	return
}

func ReviewOpenQuestV2(req []request.ReviewOpenQuestRequestV2) (err error) {
	// 开启事务
	db := global.DB.Begin()
	// 用户开放题
	userOpenQuestTimeMap := make(map[uint]time.Time)
	// 题目
	questMap := make(map[string]model.Quest)
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
		answerRes, err := sjson.Set(string(answer), fmt.Sprintf("%d", *r.Index), r.Answer)
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
		var openQuestReviewTime time.Time    // 审核时间
		var pass bool                        // 是否通过
		var userReturnScore, userScore int64 // 分数
		if openQuestReviewStatus == 2 {
			openQuestReviewTime = time.Now()
		}
		// 判断是否通过
		if openQuestReviewStatus == 2 {
			_, _, userReturnScore, userScore, pass, err = AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(answerRes), quest)
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
			UserScore:             userScore,
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
				// 创建证书
				go func() {
					GenerateCardInfo(userOpenQuest.Address, userReturnScore, request.GenerateCardInfoRequest{
						TokenId: userOpenQuest.TokenId,
						Answer:  answerRes,
					})
				}()
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

func GetUserNameTagsByAddress(address string) (nickname, name string, tags []string, err error) {
	var user model.Users
	err = global.DB.Model(&model.Users{}).
		Where("address = ?", address).
		First(&user).Error
	if err != nil {
		return
	}
	if user.Name != nil {
		name = *user.Name
	}
	if user.NickName != nil {
		nickname = *user.NickName
	}
	err = global.DB.Model(&model.UsersTag{}).Select("tag.name").
		Joins("join tag on users_tag.tag_id = tag.id").
		Where("users_tag.user_id = ?", user.ID).
		Find(&tags).Error
	return
}
