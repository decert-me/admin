package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/utils"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func AnswerCheck(key string, answer datatypes.JSON, quest model.Quest) (userReturnScore int64, pass bool, err error) {
	defer func() {
		if err != nil {
			global.LOG.Error("AnswerCheck error", zap.Error(err))
		}
	}()
	res := string(quest.MetaData)
	questData := string(quest.QuestData)
	version := gjson.Get(res, "version").Float()

	answerU, scoreList, answerS, passingScore := utils.GetAnswers(version, key, res, questData, string(answer))
	var totalScore int64
	for _, s := range scoreList {
		totalScore += s.Int()
	}
	// 获取多语言答案列表
	answers, err := GetQuestAnswersByTokenId(quest.TokenId)
	if err != nil {
		global.LOG.Error("GetQuestAnswersByTokenId error", zap.Error(err))
		return userReturnScore, false, err
	}
	// 解密答案
	var answersList [][]gjson.Result
	for _, v := range answers {
		temp := gjson.Get(utils.AnswerDecode(key, v), "@this").Array()
		answersList = append(answersList, temp) // 标准答案
		if len(answerU) != len(temp) {
			global.LOG.Error("答案数量不相等")
			return userReturnScore, false, errors.New("unexpect error")
		}
	}
	if len(answerU) != len(answerS) || len(scoreList) != len(answerS) {
		global.LOG.Error("答案数量不相等")
		return userReturnScore, false, errors.New("unexpect error")
	}
	var score int64
	for i, v := range answerS {
		if v.String() == "" {
			continue
		}
		questType := gjson.Get(v.String(), "type").String()
		questValue := gjson.Get(v.String(), "value").String()
		// 编程题目
		if questType == "coding" || questType == "special_judge_coding" {
			// 跳过不正确
			if gjson.Get(v.String(), "correct").Bool() == true {
				score += scoreList[i].Int()
			}
			continue
		}
		// 单选题
		if questType == "multiple_choice" {
			if questValue == answerU[i].String() {
				score += scoreList[i].Int()
			}
			continue
		}
		// 填空题
		if questType == "fill_blank" {
			for _, item := range answersList {
				if questValue == item[i].String() {
					score += scoreList[i].Int()
					break
				}
			}
			continue
		}
		// 多选题
		if questType == "multiple_response" {
			answerArray := gjson.Get(questValue, "@this").Array()
			fmt.Println(len(answerArray))
			fmt.Println(len(answerU[i].Array()))
			// 数量
			if len(answerArray) != len(answerU[i].Array()) {
				continue
			}
			// 内容
			allRight := true
			for _, v := range answerArray {
				var right bool
				for _, item := range answerU[i].Array() {
					if item.String() == v.String() {
						right = true
						break
					}
				}
				if !right {
					allRight = false
					break
				}
			}
			if allRight {
				score += scoreList[i].Int()
			}
		}
		if questType == "open_quest" {
			if gjson.Get(v.String(), "score").Int() != 0 {
				score += gjson.Get(v.String(), "score").Int()
			} else if gjson.Get(v.String(), "correct").Bool() == true {
				score += scoreList[i].Int()
			}
		}
	}
	if score >= passingScore {
		return score * 100 / totalScore, true, nil
	} else {
		return score * 100 / totalScore, false, nil
	}
	return
}

// IsOpenQuest 判断是否开放题
func IsOpenQuest(answerUser string) bool {
	answerU := gjson.Get(answerUser, "@this").Array()
	for _, v := range answerU {
		if v.Get("type").String() == "open_quest" {
			return true
		}
	}
	return false
}

// GetQuestAnswersByTokenId 获取题目答案
func GetQuestAnswersByTokenId(tokenId string) (answers []string, err error) {
	err = global.DB.Raw(`SELECT answer AS answers
		FROM (
		SELECT  quest_data->>'answers' AS answer FROM quest WHERE token_id = ?
		UNION
		SELECT answer FROM quest_translated WHERE token_id = ? AND answer IS NOT NULL) AS combined_data
		`, tokenId, tokenId).Scan(&answers).Error
	return
}
