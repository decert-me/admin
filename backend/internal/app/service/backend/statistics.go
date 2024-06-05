package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

// GetChallengeStatistics 挑战详情统计
func GetChallengeStatistics(r request.GetChallengeStatisticsReq) (res []response.GetQuestStatisticsRes, total int64, err error) {
	limit := r.PageSize
	offset := r.PageSize * (r.Page - 1)
	var results []response.GetQuestStatisticsRes
	countSQL := "SELECT count(1) FROM ("
	selectSQL := "SELECT * FROM ("
	dataSQL := `
		SELECT
			quest.uuid,
			quest.token_id,
			quest.ID AS quest_id,
			quest.title,
			users.address,
			users.name,
			string_agg ( DISTINCT tag.NAME, ',' ) AS tag,
			CASE 
			WHEN MAX(CAST(user_challenges.claimed AS integer))=1 THEN true
			WHEN MAX(CAST(zcloak_card.id AS integer))>0 THEN true
			ELSE false
			END AS claimed,
			CASE 
			WHEN MAX(CAST(user_challenge_log.pass AS integer))=1 THEN true
			WHEN MAX(CAST(user_open_quest.pass AS integer))=1 THEN true
			ELSE false
			END AS pass
			FROM
			"user_challenge_log"
			LEFT JOIN quest ON quest.token_id = user_challenge_log.token_id
			LEFT JOIN users ON user_challenge_log.address = users.address
			LEFT JOIN users_tag ON users_tag.user_id = users.ID 
			LEFT JOIN tag ON tag.ID = users_tag.tag_id 
			LEFT JOIN user_open_quest ON user_open_quest.token_id = quest.token_id AND user_open_quest.address= users.address
			LEFT JOIN user_challenges ON quest.token_id = user_challenges.token_id AND user_challenges.address= users.address
			LEFT JOIN zcloak_card ON zcloak_card.quest_id = quest.id AND zcloak_card.address= users.address
	`
	var whereList []string
	var valueList []interface{}
	// 应用搜索条件
	if r.SearchQuest != "" {
		whereList = append(whereList, fmt.Sprintf("(quest.title LIKE ? OR quest.token_id LIKE ?)"))
		valueList = append(valueList, "%"+r.SearchQuest+"%", "%"+r.SearchQuest+"%")
	}
	if r.SearchTag != "" {
		whereList = append(whereList, fmt.Sprintf("tag.name LIKE ?"))
		valueList = append(valueList, "%"+r.SearchTag+"%")
	}
	if r.SearchAddress != "" {
		whereList = append(whereList, fmt.Sprintf("users.address LIKE ?"))
		valueList = append(valueList, "%"+r.SearchAddress+"%")
	}
	dataSQL += " WHERE quest.token_id is not null AND  users.address is not null"
	if len(whereList) > 0 {
		dataSQL += " AND " + strings.Join(whereList, " AND ")
	}
	dataSQL += " GROUP BY users.address,quest.ID,users.name) as tt"
	// 分页SQL
	paginateSQL := " LIMIT ? OFFSET ?"
	// 过滤条件
	whereList = nil // 清空
	if r.Pass != nil {
		if *r.Pass {
			whereList = append(whereList, "pass = true")
		} else {
			whereList = append(whereList, "pass =false")
		}
	}
	if r.Claimed != nil {
		if *r.Claimed {
			whereList = append(whereList, "claimed = true")
		} else {
			whereList = append(whereList, "claimed =false")
		}
	}
	if len(whereList) > 0 {
		dataSQL += " WHERE " + strings.Join(whereList, " AND ")
	}
	// 执行查询
	db := global.DB
	// 获取总数用于分页
	err = db.Raw(countSQL+dataSQL, valueList...).Scan(&total).Error
	if err != nil {
		return
	}
	// 执行查询
	valueList = append(valueList, limit, offset)
	err = db.Raw(selectSQL+dataSQL+paginateSQL, valueList...).Scopes(Paginate(r.Page, r.PageSize)).Find(&results).Error
	if err != nil {
		return res, total, err
	}
	db = global.DB
	// 处理数据
	for i, v := range results {
		// 查询题目
		var quest model.Quest
		if err := db.Where("token_id = ?", v.TokenID).First(&quest).Error; err != nil {
			return res, total, err
		}
		// 判断是否已经领取
		var userChallenge model.UserChallenges
		if err := db.Where("token_id = ? AND address = ?", v.TokenID, v.Address).First(&userChallenge).Error; err == nil {
			results[i].Claimed = true
			results[i].Pass = true
		}
		var zkCard model.ZcloakCard
		if err := db.Where("address = ? AND quest_id = ?", v.Address, v.QuestID).First(&zkCard).Error; err == nil {
			results[i].Claimed = true
			results[i].Pass = true
		}
		isOpenQuest := IsOpenQuest(gjson.Get(string(quest.QuestData), "questions").String()) // 是否开放题
		// 获取及格分数
		passingScore := gjson.Get(string(quest.QuestData), "passingScore").Int()
		// 查询挑战记录
		if isOpenQuest {
			var userOpenQuest model.UserOpenQuest
			userOpenQuestDB := global.DB
			if results[i].Claimed {
				userOpenQuestDB.Where("pass=true")
			}
			if err := userOpenQuestDB.Where("token_id = ? AND address = ?", v.TokenID, v.Address).Order("id desc").First(&userOpenQuest).Error; err == nil {
				results[i].ChallengeTime = userOpenQuest.CreatedAt
			}
			_, userReturnRawScore, _, _, _, _ := AnswerCheck(global.CONFIG.Quest.EncryptKey, userOpenQuest.Answer, quest)
			results[i].ScoreDetail = strconv.Itoa(int(userReturnRawScore)) + "/" + strconv.Itoa(int(passingScore))
			// 获取批注
			temp := gjson.Get(string(userOpenQuest.Answer), "@this").Array()
			for ii, vv := range temp {
				title := gjson.Get(string(quest.QuestData), "questions."+strconv.Itoa(ii)+".title").String()
				annotation := gjson.Get(vv.String(), "annotation").String()
				if annotation == "" {
					continue
				}
				results[i].Annotation += fmt.Sprintf("第 %d 题 %s \n %s \n", ii+1, title, annotation)
			}
		} else {
			var userChallengeLog model.UserChallengeLog
			userChallengeLogDB := global.DB
			if results[i].Claimed {
				userChallengeLogDB.Where("pass=true")
			}
			if err := userChallengeLogDB.Where("token_id = ? AND address = ?", v.TokenID, v.Address).Order("id desc").First(&userChallengeLog).Error; err == nil {
				results[i].ChallengeTime = userChallengeLog.CreatedAt
			}
			_, userReturnRawScore, _, _, _, _ := AnswerCheck(global.CONFIG.Quest.EncryptKey, userChallengeLog.Answer, quest)
			results[i].ScoreDetail = strconv.Itoa(int(userReturnRawScore)) + "/" + strconv.Itoa(int(passingScore))
		}

	}

	return results, total, nil
}

// GetChallengeUserStatistics 挑战者统计
func GetChallengeUserStatistics(r request.GetChallengeUserStatisticsReq) (res []response.GetChallengeUserStatisticsRes, total int64, err error) {
	var results []response.GetChallengeUserStatisticsRes
	db := global.DB.Table("users").
		Select("users.id as user_id, users.address, users.name,string_agg(tag.name, ',') as tags").
		Joins("LEFT JOIN users_tag ON users_tag.user_id = users.id").
		Joins("LEFT JOIN tag ON tag.id = users_tag.tag_id").
		Group("users.id")

	if r.SearchTag != "" {
		db = db.Where("tag.name = ?", r.SearchTag)
	}

	if r.SearchAddress != "" {
		db = db.Where("users.address = ?", r.SearchAddress)
	}
	// 获取总数用于分页
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	// 执行查询
	err = db.Scopes(Paginate(r.Page, r.PageSize)).Find(&results).Error
	if err != nil {
		return res, total, err
	}
	db = global.DB
	// 统计数据
	for i, v := range results {
		// 领取NFT数量
		db.Raw(`SELECT COUNT(1) FROM
			(
			SELECT
				quest.ID 
			FROM
				"user_challenges"
				LEFT JOIN quest ON quest.token_id = user_challenges.token_id 
			WHERE
				address = ? 
			GROUP BY
				quest.ID 
			UNION
			SELECT
				quest_id 
			FROM
				zcloak_card 
			WHERE
				address = ? 
			GROUP BY
			quest_id ) AS f`, v.Address, v.Address).Scan(&results[i].ClaimNum)
		// 挑战成功/失败数量
		type CountResult struct {
			TokenId      string `json:"token_id"`
			PassCount    int    `json:"pass_count"`
			NotPassCount int    `json:"not_pass_count"`
		}
		var countResult []CountResult
		if err := global.DB.Raw(`
		SELECT 
			token_id,
			sum(pass_count) as pass_count,
			sum(not_pass_count) as not_pass_count
		FROM (
			(SELECT token_id, sum(case when pass then 1 else 0 end) as pass_count, sum(case when not pass then 1 else 0 end) as not_pass_count 
			 FROM user_challenge_log
			 WHERE address = ?
			 GROUP BY token_id)
			UNION ALL
			(SELECT token_id, sum(case when pass then 1 else 0 end) as pass_count, sum(case when not pass then 1 else 0 end) as not_pass_count 
			 FROM user_open_quest
			 WHERE address = ?
			 GROUP BY token_id)
		) as combined
		GROUP BY token_id
		`, v.Address, v.Address).Scan(&countResult).Error; err != nil {
			return res, total, err
		}
		for _, result := range countResult {
			if result.PassCount == 0 {
				results[i].FailNum += 1
				continue
			}
			results[i].SuccessNum += 1
		}
		results[i].NotClaimNum = results[i].SuccessNum - results[i].ClaimNum
	}
	return results, total, nil
}
