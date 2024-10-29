package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	result, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, r.Answer, quest)
	score := result.UserReturnScore
	// 写入审核结果
	err = global.DB.Model(&model.UserOpenQuest{}).Where("id = ? AND open_quest_review_status = 1", r.ID).Updates(&model.UserOpenQuest{
		OpenQuestReviewTime:   time.Now(),
		OpenQuestReviewStatus: 2,
		OpenQuestScore:        score,
		Answer:                r.Answer,
		Pass:                  result.Pass,
		UserScore:             result.UserScore,
	}).Error
	// 写入Message
	var message model.UserMessage
	if result.Pass {
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
	countSQL := `
		SELECT
			count(1)
		FROM
			quest,
			jsonb_array_elements(quest.quest_data -> 'questions') WITH ORDINALITY AS t(json_element, idx)
		WHERE
			t.json_element ->> 'type' = 'open_quest' 
			AND quest.status = 1
	`
	err = db.Raw(countSQL).Scan(&total).Error
	dataSQL := `
		WITH quest_data AS (
			SELECT
				t.json_element ->> 'title' AS title,
				quest.title AS challenge_title,
				quest.uuid,
				quest.token_id,
				(t.idx::int - 1) AS index,
				quest.add_ts AS add_ts,
				quest.status
			FROM
				quest,
				jsonb_array_elements(quest.quest_data -> 'questions') WITH ORDINALITY AS t(json_element, idx)
			WHERE
				t.json_element ->> 'type' = 'open_quest' 
				AND quest.status = 1
		)
		SELECT 
			quest_data.title, 
			quest_data.challenge_title, 
			quest_data.uuid, 
			quest_data.token_id, 
			quest_data.index, 
			quest_data.add_ts,
			COUNT(user_open_quest.token_id) AS to_review_count
		FROM
			quest_data
		LEFT JOIN user_open_quest ON quest_data.token_id = user_open_quest.token_id
			AND user_open_quest.deleted_at IS NULL
			AND quest_data.status = 1
			AND EXISTS (
				SELECT 1
				FROM jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx)
				WHERE t.json_element ->> 'type' = 'open_quest'
				AND idx = index+1
				AND t.json_element ->> 'score' IS NULL
				AND t.json_element ->> 'correct' IS NULL
			)
		GROUP BY 
			quest_data.title, 
			quest_data.challenge_title, 
			quest_data.uuid, 
			quest_data.token_id, 
			quest_data.index, 
			quest_data.add_ts
		ORDER BY to_review_count desc
		LIMIT ? OFFSET ? 
	`
	err = db.Raw(dataSQL, limit, offset).Scan(&list).Error
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
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
	}
	return list, total, totalToReview, nil
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
					user_open_quest.updated_at as updated_at,
					user_open_quest.created_at as created_at,
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
		// 提交次数
		submitCountSQL := `
		SELECT 
			count(1)
		FROM
			user_open_quest
		WHERE
			user_open_quest.deleted_at IS NULL AND user_open_quest.token_id = ? AND user_open_quest.address = ? AND user_open_quest.id <= ?
		`
		submitCountErr := global.DB.Raw(submitCountSQL, list[i].TokenId, list[i].Address, list[i].ID).Scan(&list[i].SubmitCount).Error
		if submitCountErr != nil {
			if submitCountErr == gorm.ErrRecordNotFound {
				list[i].SubmitCount = 0
			} else {
				log.Error("获取提交次数失败", "error", submitCountErr)
			}
		}
		// 上次分数
		lastAnswerSQL := `
		SELECT
			COALESCE(json_element->>'score', '0') AS score
		FROM
			user_open_quest
		JOIN
			jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
		WHERE
			user_open_quest.deleted_at IS NULL AND json_element->>'type' = 'open_quest' AND user_open_quest.token_id = ? AND idx = ? AND user_open_quest.address = ? AND open_quest_review_status=2 AND user_open_quest.id < ?
		ORDER BY
			updated_at DESC
		LIMIT 1
		`
		lastAnswerErr := global.DB.Raw(lastAnswerSQL, list[i].TokenId, list[i].Index+1, list[i].Address, list[i].ID).Scan(&list[i].LastScore).Error
		if lastAnswerErr != nil {
			if lastAnswerErr == gorm.ErrRecordNotFound {
				list[i].LastScore = 0
			} else {
				log.Error("获取上次分数失败", "error", lastAnswerErr)
			}
		}
		quest := model.Quest{
			TokenId:   list[i].TokenId,
			MetaData:  list[i].MetaData,
			QuestData: list[i].QuestData,
		}
		result, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, list[i].UserAnswer, quest)
		if err != nil {
			return list, 0, err
		}
		list[i].TotalScore = result.TotalScore
		list[i].UserScore = result.UserScore
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
			result, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(answerRes), quest)
			if err != nil {
				db.Rollback()
				return errors.New("服务器错误")
			}
			userReturnScore = result.UserReturnScore
			userScore = result.UserScore
			pass = result.Pass
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

func GetUserQuestDetail(r request.GetUserQuestDetailRequest) (res response.GetUserQuestDetailResponse, err error) {
	address := r.Address
	res.Address = address
	err = global.DB.Model(&model.Quest{}).
		Select("quest.*,COALESCE(tr.title,quest.title) as title,COALESCE(tr.description,quest.description) as description,"+
			"COALESCE(tr.meta_data,quest.meta_data) as meta_data,COALESCE(tr.quest_data,quest.quest_data) as quest_data,"+
			"b.claimed,b.user_score,b.nft_address,b.badge_token_id,b.chain_id as badge_chain_id,COALESCE(o.open_quest_review_status,0) as open_quest_review_status,COALESCE(o.answer,l.answer) as answer,COALESCE(o.created_at,l.created_at) as submit_time").
		Joins("left join user_challenges b ON quest.token_id=b.token_id AND b.address= ?", address).
		Joins("left join user_challenge_log l ON quest.token_id=l.token_id AND l.address= ? AND l.deleted_at IS NULL", address).
		Joins("left join user_open_quest o ON quest.token_id=o.token_id AND o.address= ? AND o.deleted_at IS NULL", address).
		Joins("LEFT JOIN quest_translated tr ON quest.token_id = tr.token_id AND tr.language = ?", "zh-CN").
		Where("quest.uuid", r.UUID).
		Order("o.pass desc,l.pass desc,l.add_ts desc,o.id desc").
		First(&res).Error
	if err != nil {
		return res, err
	}
	// 获取所有答案
	err = global.DB.Raw(`SELECT answer AS answers
		FROM (
		SELECT  quest_data->>'answers' AS answer FROM quest WHERE token_id = ?
		UNION
		SELECT answer FROM quest_translated WHERE token_id = ? AND answer IS NOT NULL) AS combined_data
		`, res.TokenId, res.TokenId).Scan(&res.Answers).Error

	var quest model.Quest
	if err = global.DB.Model(&model.Quest{}).Where("token_id = ?", res.TokenId).First(&quest).Error; err != nil {
		return
	}
	result, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, res.Answer, quest)
	if err != nil {
		return response.GetUserQuestDetailResponse{}, err
	}
	answer := gjson.Get(string(res.Answer), "@this").Array()
	answerStr := gjson.Get(string(res.Answer), "@this").String()
	for i, _ := range answer {
		answerRes, err := sjson.Set(answerStr, fmt.Sprintf("%d.score", i), result.UserScoreList[i])
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		answerStr = answerRes
	}
	res.Answer = datatypes.JSON(answerStr)
	return
}
