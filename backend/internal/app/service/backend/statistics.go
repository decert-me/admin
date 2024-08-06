package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
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
			string_agg ( DISTINCT tag.NAME, ',' ) AS tags,
			CASE 
			WHEN MAX(CAST(user_challenges.claimed AS integer))=1 THEN true
			WHEN MAX(CAST(zcloak_card.id AS integer))>0 THEN true
			ELSE false
			END AS claimed,
			CASE 
			WHEN MAX(CAST(user_challenge_log.pass AS integer))=1 THEN true
			WHEN MAX(CAST(user_open_quest.pass AS integer))=1 THEN true
			ELSE false
			END AS pass,
			CASE 
			WHEN MAX(CASE WHEN user_open_quest.open_quest_review_status = 1 THEN 1 ELSE 0 END) = 1 THEN true
			ELSE false
			END AS reviewing
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
		whereList = append(whereList, fmt.Sprintf("(quest.title LIKE ? OR quest.token_id LIKE ? OR quest.uuid LIKE ?)"))
		valueList = append(valueList, "%"+r.SearchQuest+"%", "%"+r.SearchQuest+"%", "%"+r.SearchQuest+"%")
	}
	if r.SearchTag != "" {
		whereList = append(whereList, fmt.Sprintf("tag.name LIKE ?"))
		valueList = append(valueList, "%"+r.SearchTag+"%")
	}
	if r.SearchAddress != "" {
		whereList = append(whereList, fmt.Sprintf("(users.address ILIKE ? OR users.name ILIKE ?)"))
		valueList = append(valueList, "%"+r.SearchAddress+"%", "%"+r.SearchAddress+"%")
	}
	dataSQL += " WHERE quest.token_id is not null AND  users.address is not null AND quest.disabled = false"
	if len(whereList) > 0 {
		dataSQL += " AND " + strings.Join(whereList, " AND ")
	}
	dataSQL += " GROUP BY users.address,quest.ID,users.name) as tt"
	// 分页SQL
	paginateSQL := " LIMIT ? OFFSET ?"
	if r.SearchTag != "" && r.SearchQuest != "" {
		paginateSQL = ""
	}
	// 过滤条件
	whereList = nil // 清空
	if r.Pass != nil {
		if *r.Pass {
			whereList = append(whereList, "pass = true")
		} else {
			whereList = append(whereList, "pass = false")
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
	if !(r.SearchTag != "" && r.SearchQuest != "") {
		valueList = append(valueList, limit, offset)
	}

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
		var userChallengeLog model.UserChallengeLog
		userChallengeLogDB := global.DB
		if results[i].Claimed {
			userChallengeLogDB.Where("pass = true")
		}
		if err := userChallengeLogDB.Where("token_id = ? AND address = ?", v.TokenID, v.Address).Order("pass desc,id desc").First(&userChallengeLog).Error; err == nil {
			results[i].ChallengeTime = userChallengeLog.CreatedAt
		}
		_, userReturnRawScore, _, _, _, _ := AnswerCheck(global.CONFIG.Quest.EncryptKey, userChallengeLog.Answer, quest)
		results[i].ScoreDetail = strconv.Itoa(int(userReturnRawScore)) + "/" + strconv.Itoa(int(passingScore))
		// 开放题
		if isOpenQuest {
			var userOpenQuest model.UserOpenQuest
			userOpenQuestDB := global.DB
			if results[i].Claimed {
				userOpenQuestDB.Where("pass = true")
			}
			if err := userOpenQuestDB.Where("token_id = ? AND address = ?", v.TokenID, v.Address).Order("pass desc,id desc").First(&userOpenQuest).Error; err == nil {
				results[i].ChallengeTime = userOpenQuest.CreatedAt
			} else {
				continue
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
		}
		if results[i].Claimed || results[i].Pass {
			results[i].ChallengeResult = "成功"
		} else if results[i].Reviewing {
			results[i].ChallengeResult = "评分中"
		} else {
			results[i].ChallengeResult = "失败"
		}
	}

	// 额外处理
	if r.SearchTag != "" && r.SearchQuest != "" {
		// 查询所有挑战
		var questList []model.Quest
		err = global.DB.Model(&model.Quest{}).Where("quest.title LIKE ? OR quest.token_id LIKE ? OR quest.uuid LIKE ?", r.SearchQuest, r.SearchQuest, r.SearchQuest).Find(&questList).Error
		if err != nil {
			log.Error("Failed to query quest", zap.Error(err))
			return results, total, err
		}
		// 查询所有地址
		var userList []model.Users
		err = global.DB.Model(&model.Users{}).
			Joins("LEFT JOIN users_tag ON users_tag.user_id = users.ID").
			Joins("LEFT JOIN tag ON tag.id = users_tag.tag_id").
			Where("tag.name LIKE ?", "%"+r.SearchTag+"%").
			Find(&userList).Error

		for _, user := range userList {
			for _, quest := range questList {
				if quest.Disabled {
					continue
				}
				if isAddressAndTokenIDInResults(user.Address, quest.TokenId, results) {
					continue
				}
				var result response.GetQuestStatisticsRes
				result.Address = user.Address
				result.TokenID = quest.TokenId
				result.QuestID = int64(quest.ID)
				result.Title = quest.Title
				result.Name = *user.Name
				result.Tags = r.SearchTag
				result.ChallengeTime = time.Time{}
				result.Pass = false
				result.Claimed = false
				result.ChallengeResult = "-"
				results = append(results, result)
				total++
			}
		}
	}
	// 手动分页
	if len(results) > int(r.PageSize) {
		results = results[offset : offset+limit]
	}
	return results, total, nil
}

func isAddressAndTokenIDInResults(address string, tokenID string, results []response.GetQuestStatisticsRes) bool {
	for _, result := range results {
		if result.Address == address && result.TokenID == tokenID {
			return true
		}
	}
	return false
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
		db = db.Where("tag.name like ?", "%"+r.SearchTag+"%")
	}

	if r.SearchAddress != "" {
		db = db.Where("(users.address ILIKE ? OR users.name ILIKE ?)", "%"+r.SearchAddress+"%", "%"+r.SearchAddress+"%")
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
			LEFT JOIN user_challenge_log ON user_challenge_log.token_id = user_challenges.token_id AND user_challenge_log.address = user_challenges.address
			WHERE
				user_challenges.address = ? AND quest.token_id IS NOT NULL AND user_challenge_log.token_id IS NOT NULL 
			GROUP BY
				quest.ID 
			UNION
			SELECT
				quest_id 
			FROM
				zcloak_card 
			LEFT JOIN quest ON quest.id = zcloak_card.quest_id 
						LEFT JOIN user_challenge_log ON user_challenge_log.token_id = quest.token_id AND user_challenge_log.address = zcloak_card.address
			WHERE
				zcloak_card.address = ? AND quest.token_id IS NOT NULL AND user_challenge_log.token_id IS NOT NULL 
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
			(SELECT user_challenge_log.token_id, sum(case when pass then 1 else 0 end) as pass_count, sum(case when not pass then 1 else 0 end) as not_pass_count 
			 FROM user_challenge_log
			 LEFT JOIN quest ON user_challenge_log.token_id=quest.token_id
			 WHERE address = ? AND quest.token_id IS NOT NULL
			 GROUP BY user_challenge_log.token_id)
			UNION ALL
			(SELECT user_open_quest.token_id, sum(case when pass then 1 else 0 end) as pass_count, sum(case when not pass then 1 else 0 end) as not_pass_count 
			 FROM user_open_quest
			 LEFT JOIN quest ON user_open_quest.token_id=quest.token_id
			 WHERE address = ? AND quest.token_id IS NOT NULL
			 GROUP BY user_open_quest.token_id)
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
		if results[i].NotClaimNum < 0 {
			results[i].NotClaimNum = 0
		}
	}
	return results, total, nil
}

// GetChallengeStatisticsSummary 挑战详情总计
func GetChallengeStatisticsSummary(r request.GetChallengeStatisticsReq) (res response.GetQuestStatisticsSummaryRes, err error) {
	selectSQL := `
	SELECT  COUNT(DISTINCT token_id) AS challenge_num,
    		COUNT(DISTINCT address) AS challenge_user_num,
			SUM(CASE WHEN pass = true THEN 1 ELSE 0 END) AS success_num,
    		SUM(CASE WHEN pass = false THEN 1 ELSE 0 END) AS fail_num,
			SUM(CASE WHEN claimed = true THEN 1 ELSE 0 END) AS claim_num,
			SUM(CASE WHEN claimed = false AND pass = true THEN 1 ELSE 0 END) AS not_claim_num
		FROM (`
	dataSQL := `
		SELECT
			quest.uuid,
			quest.token_id,
			quest.ID AS quest_id,
			quest.title,
			users.address,
			users.name,
			string_agg ( DISTINCT tag.NAME, ',' ) AS tags,
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
			FROM "user_challenge_log"
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
		whereList = append(whereList, fmt.Sprintf("(quest.title LIKE ? OR quest.token_id LIKE ? OR quest.uuid LIKE ?)"))
		valueList = append(valueList, "%"+r.SearchQuest+"%", "%"+r.SearchQuest+"%", "%"+r.SearchQuest+"%")
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
	// 执行查询
	err = db.Raw(selectSQL+dataSQL, valueList...).Find(&res).Error
	if err != nil {
		return res, err
	}
	return
}
