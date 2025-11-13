package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/backend"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

// GetAiJudgeConfigList 获取AI判题配置列表
func GetAiJudgeConfigList(c *gin.Context) {
	var configs []model.AiJudgeConfig

	if err := global.DB.Order("created_at DESC").Find(&configs).Error; err != nil {
		global.LOG.Error("获取AI判题配置列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(configs, c)
}

// CreateAiJudgeConfig 创建AI判题配置
func CreateAiJudgeConfig(c *gin.Context) {
	var req request.CreateAiJudgeConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	config := model.AiJudgeConfig{
		Title:   req.Title,
		Config:  req.Config,
		Enabled: false,
	}

	if err := global.DB.Create(&config).Error; err != nil {
		global.LOG.Error("创建AI判题配置失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}

	response.OkWithData(config, c)
}

// UpdateAiJudgeConfig 更新AI判题配置
func UpdateAiJudgeConfig(c *gin.Context) {
	var req request.UpdateAiJudgeConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	var config model.AiJudgeConfig
	if err := global.DB.First(&config, req.ID).Error; err != nil {
		response.FailWithMessage("配置不存在", c)
		return
	}

	config.Title = req.Title
	config.Config = req.Config

	if err := global.DB.Save(&config).Error; err != nil {
		global.LOG.Error("更新AI判题配置失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithData(config, c)
}

// DeleteAiJudgeConfig 删除AI判题配置
func DeleteAiJudgeConfig(c *gin.Context) {
	var req request.DeleteAiJudgeConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	if err := global.DB.Delete(&model.AiJudgeConfig{}, req.ID).Error; err != nil {
		global.LOG.Error("删除AI判题配置失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// ToggleAiJudgeConfig 切换AI判题配置启用状态（确保只有一个配置被启用）
func ToggleAiJudgeConfig(c *gin.Context) {
	var req request.ToggleAiJudgeConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	var config model.AiJudgeConfig
	if err := global.DB.First(&config, req.ID).Error; err != nil {
		response.FailWithMessage("配置不存在", c)
		return
	}

	// 开始事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 如果要启用这个配置，先禁用所有其他配置
	if !config.Enabled {
		if err := tx.Model(&model.AiJudgeConfig{}).Where("enabled = ?", true).Update("enabled", false).Error; err != nil {
			tx.Rollback()
			global.LOG.Error("禁用其他配置失败", zap.Error(err))
			response.FailWithMessage("操作失败", c)
			return
		}
		config.Enabled = true
	} else {
		// 如果要禁用当前配置
		config.Enabled = false
	}

	if err := tx.Save(&config).Error; err != nil {
		tx.Rollback()
		global.LOG.Error("更新配置状态失败", zap.Error(err))
		response.FailWithMessage("操作失败", c)
		return
	}

	tx.Commit()
	response.OkWithData(config, c)
}

// GetEnabledAiJudgeConfig 获取当前启用的AI判题配置
func GetEnabledAiJudgeConfig(c *gin.Context) {
	var config model.AiJudgeConfig

	if err := global.DB.Where("enabled = ?", true).First(&config).Error; err != nil {
		response.FailWithMessage("当前没有启用的AI配置", c)
		return
	}

	response.OkWithData(config, c)
}

// AIConfig AI配置结构体（用于解析存储在数据库中的配置）
type AIConfig struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
}

// AiGrade AI判题接口
func AiGrade(c *gin.Context) {
	var req request.AiGradeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	// 获取启用的AI配置
	var config model.AiJudgeConfig
	if err := global.DB.Where("enabled = ?", true).First(&config).Error; err != nil {
		response.FailWithMessage("当前没有启用的AI配置，请先在AI配置中启用一个配置", c)
		return
	}

	// 解析配置内容（从富文本中提取JSON）
	aiConfig, err := parseAIConfig(config.Config)
	if err != nil {
		global.LOG.Error("解析AI配置失败", zap.Error(err))
		response.FailWithMessage("AI配置格式错误，请检查配置内容", c)
		return
	}

	// 构建提示词
	systemPrompt := fmt.Sprintf(
		"你是一名web3开发专家，现在有一些web3开发学习过程中的题目和答案，需要你来进行批改，本题总分%d分，及格分%d分，请根据用户的答案和题目的要求来进行批改。如果通过，那么直接返回用户得分多少即可，如果不通过，那么需要返回不通过的分数和不通过的理由作为批注。如果答案中存在链接，那么需要访问链接中的具体网页来获取信息进行判题。如果存在附件，需要查看附件中的具体内容来进行判题。",
		req.QuestionScore,
		req.PassScore,
	)

	// 构建用户提示词，包含附件信息
	userPrompt := fmt.Sprintf("题目：%s\n\n用户答案：%s", req.QuestionTitle, req.UserAnswer)

	// 如果有附件，添加附件信息
	if len(req.AttachmentUrls) > 0 {
		userPrompt += "\n\n用户提交的附件："
		for i, url := range req.AttachmentUrls {
			userPrompt += fmt.Sprintf("\n%d. %s", i+1, url)
		}
	}

	userPrompt += "\n\n请根据以上信息进行判题，返回格式：得分：X分\n批注：（如果不通过则说明理由）"

	// 调用AI API
	result, err := callAIAPI(aiConfig, systemPrompt, userPrompt)
	if err != nil {
		global.LOG.Error("调用AI API失败", zap.Error(err))
		response.FailWithMessage("AI判题失败："+err.Error(), c)
		return
	}

	// 解析AI返回结果，提取分数和批注
	score, annotation := parseAIResult(result, req.QuestionScore)

	// 判断是否通过：分数达到及格分则通过
	isPassed := score >= req.PassScore

	// 如果通过，清空批注；如果不通过，保留批注
	finalAnnotation := ""
	if !isPassed {
		finalAnnotation = annotation
	}

	response.OkWithData(gin.H{
		"score":         score,
		"annotation":    finalAnnotation,
		"raw_result":    result,
		"system_prompt": systemPrompt,
		"user_prompt":   userPrompt,
		"is_passed":     isPassed,
	}, c)
}

// parseAIConfig 解析存储在数据库中的AI配置（从富文本HTML中提取JSON）
func parseAIConfig(configHTML string) (*AIConfig, error) {
	// 移除HTML标签，提取纯文本
	re := regexp.MustCompile(`<[^>]*>`)
	configText := re.ReplaceAllString(configHTML, "")

	// 移除HTML实体
	configText = strings.ReplaceAll(configText, "&nbsp;", " ")
	configText = strings.ReplaceAll(configText, "&quot;", "\"")
	configText = strings.ReplaceAll(configText, "&lt;", "<")
	configText = strings.ReplaceAll(configText, "&gt;", ">")
	configText = strings.ReplaceAll(configText, "&amp;", "&")

	// 去除多余空白字符
	configText = strings.TrimSpace(configText)

	// 尝试解析JSON
	var aiConfig AIConfig
	if err := json.Unmarshal([]byte(configText), &aiConfig); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	// 验证必填字段
	if aiConfig.APIKey == "" || aiConfig.BaseURL == "" {
		return nil, fmt.Errorf("配置缺少必填字段 apiKey 或 baseUrl")
	}

	// 设置默认模型
	if aiConfig.Model == "" {
		aiConfig.Model = "gpt-3.5-turbo"
	}

	return &aiConfig, nil
}

// callAIAPI 调用AI API
func callAIAPI(config *AIConfig, systemPrompt, userPrompt string) (string, error) {
	// 构建请求体
	requestBody := map[string]interface{}{
		"model": config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %w", err)
	}

	// 创建HTTP请求
	apiURL := strings.TrimRight(config.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取AI回复内容
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("响应格式错误：缺少choices字段")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误：choices格式不正确")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误：缺少message字段")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("响应格式错误：content不是字符串")
	}

	return content, nil
}

// parseAIResult 解析AI返回结果，提取分数和批注
func parseAIResult(result string, maxScore int) (int, string) {
	// 第一步：预处理 - 去除所有的星号"*"
	result = strings.ReplaceAll(result, "*", "")

	// 默认值
	score := 0
	annotation := result

	// 尝试提取分数（多种格式）
	// 格式1: "得分：X分" 或 "得分: X分"
	scorePattern1 := regexp.MustCompile(`得分[：:]\s*(\d+)\s*分`)
	if matches := scorePattern1.FindStringSubmatch(result); len(matches) > 1 {
		if s, err := strconv.Atoi(matches[1]); err == nil {
			score = s
		}
	}

	// 格式2: "X分" (单独出现的数字+分)
	if score == 0 {
		scorePattern2 := regexp.MustCompile(`(\d+)\s*分`)
		if matches := scorePattern2.FindStringSubmatch(result); len(matches) > 1 {
			if s, err := strconv.Atoi(matches[1]); err == nil {
				score = s
			}
		}
	}

	// 格式3: "score: X" 或 "分数: X"
	if score == 0 {
		scorePattern3 := regexp.MustCompile(`(?:score|分数)[：:]\s*(\d+)`)
		if matches := scorePattern3.FindStringSubmatch(result); len(matches) > 1 {
			if s, err := strconv.Atoi(matches[1]); err == nil {
				score = s
			}
		}
	}

	// 限制分数不超过最大值
	if score > maxScore {
		score = maxScore
	}

	// 尝试提取批注
	// 格式1: "批注：xxx" 或 "批注: xxx"  (支持多行内容)
	annotationPattern1 := regexp.MustCompile(`批注[：:]\s*([\s\S]+?)(?:\n\n|$)`)
	if matches := annotationPattern1.FindStringSubmatch(result); len(matches) > 1 {
		annotation = strings.TrimSpace(matches[1])
	} else {
		// 如果没有明确的批注标记，尝试提取"得分"之后的所有内容
		scoreLinePattern := regexp.MustCompile(`得分[：:]\s*\d+\s*分\s*(.*)`)
		if matches := scoreLinePattern.FindStringSubmatch(result); len(matches) > 1 {
			annotation = strings.TrimSpace(matches[1])
		} else {
			// 如果都没有，使用整个结果作为批注
			annotation = result
		}
	}

	return score, annotation
}

// ToggleAutoGrading 切换自动判题状态
func ToggleAutoGrading(c *gin.Context) {
	var req request.ToggleAutoGradingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	var config model.AiJudgeConfig
	if err := global.DB.First(&config, req.ID).Error; err != nil {
		response.FailWithMessage("配置不存在", c)
		return
	}

	// 更新自动判题状态
	if err := global.DB.Model(&config).Update("auto_grading", req.AutoGrading).Error; err != nil {
		global.LOG.Error("更新自动判题状态失败", zap.Error(err))
		response.FailWithMessage("操作失败", c)
		return
	}

	response.OkWithMessage("操作成功", c)
}

// GetPendingGradeList 获取待判题列表
func GetPendingGradeList(c *gin.Context) {
	var total int64

	// 使用与评分列表完全相同的逻辑：统计所有未评分的开放题答案（而不是记录数）
	// 这里统计的是答案数量，一条 user_open_quest 记录可能包含多个开放题答案
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
			user_open_quest.deleted_at IS NULL
			AND quest.status = 1
			AND json_element->>'type' = 'open_quest'
			AND json_element->>'score' IS NULL
			AND json_element->>'correct' IS NULL
	`

	if err := global.DB.Raw(countSQL).Scan(&total).Error; err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(gin.H{
		"list":  []interface{}{},
		"total": total,
	}, c)
}

// BatchGrade 批量AI判题
func BatchGrade(c *gin.Context) {
	var req request.BatchGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	maxCount := req.MaxCount
	if maxCount <= 0 || maxCount > 50 {
		maxCount = 10 // 默认10条，最多50条
	}

	// 获取启用的AI配置
	var config model.AiJudgeConfig
	if err := global.DB.Where("enabled = ?", true).First(&config).Error; err != nil {
		response.FailWithMessage("当前没有启用的AI配置，请先在AI配置中启用一个配置", c)
		return
	}

	// 获取待判题列表：需要记录ID、token_id和答案索引
	type PendingQuest struct {
		ID      uint   `json:"id"`
		TokenId string `json:"token_id"`
		Index   int    `json:"index"` // 答案数组中的索引（从0开始）
	}
	var pendingQuests []PendingQuest

	querySQL := `
		SELECT
			user_open_quest.id,
			user_open_quest.token_id,
			(t.idx::int - 1) AS index
		FROM
			user_open_quest
		JOIN
			jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
		JOIN
			quest ON quest.token_id = user_open_quest.token_id
		WHERE
			user_open_quest.deleted_at IS NULL
			AND quest.status = 1
			AND json_element->>'type' = 'open_quest'
			AND json_element->>'score' IS NULL
			AND json_element->>'correct' IS NULL
		ORDER BY user_open_quest.id ASC
		LIMIT ?
	`

	if err := global.DB.Raw(querySQL, maxCount).Scan(&pendingQuests).Error; err != nil {
		response.FailWithMessage("获取待判题列表失败", c)
		return
	}

	if len(pendingQuests) == 0 {
		response.OkWithMessage("没有待判题的开放题", c)
		return
	}

	// 异步批量判题
	go func() {
		for _, pq := range pendingQuests {
			ProcessOneOpenQuest(pq.ID, pq.TokenId, pq.Index, config)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	response.OkWithData(gin.H{
		"message": "开始批量判题",
		"count":   len(pendingQuests),
	}, c)
}

// ProcessOneOpenQuest 处理单个开放题答案的AI判题
// recordID: user_open_quest 记录ID
// tokenId: 题目token_id
// answerIndex: 答案数组中的索引位置（从0开始）
// config: AI配置
func ProcessOneOpenQuest(recordID uint, tokenId string, answerIndex int, config model.AiJudgeConfig) {
	// 1. 获取记录（加锁）
	var userOpenQuest model.UserOpenQuest
	// 使用FOR UPDATE锁定记录，防止并发修改
	if err := global.DB.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", recordID).
		First(&userOpenQuest).Error; err != nil {
		global.LOG.Error("获取记录失败", zap.Error(err), zap.Uint("record_id", recordID))
		return
	}

	// 2. 检查该答案是否已经评分
	answerJSON := string(userOpenQuest.Answer)
	answerPath := fmt.Sprintf("%d", answerIndex)

	// 检查是否已有score或correct字段
	if gjson.Get(answerJSON, answerPath+".score").Exists() ||
	   gjson.Get(answerJSON, answerPath+".correct").Exists() {
		// 已评分，跳过
		global.LOG.Info("答案已评分，跳过", zap.Uint("record_id", recordID), zap.Int("index", answerIndex))
		return
	}

	// 3. 获取题目信息
	var questInfo model.Quest
	if err := global.DB.Where("token_id = ?", tokenId).First(&questInfo).Error; err != nil {
		global.LOG.Error("获取题目失败", zap.Error(err), zap.String("token_id", tokenId))
		return
	}

	// 4. 提取该答案的内容
	answerData := gjson.Get(answerJSON, answerPath)
	answerValue := answerData.Get("value").String()

	// 5. 提取附件URL
	var attachmentUrls []string
	annexArray := answerData.Get("annex").Array()
	for _, annex := range annexArray {
		hash := annex.Get("hash").String()
		name := annex.Get("name").String()
		if hash != "" {
			url := fmt.Sprintf("https://ipfs.decert.me/ipfs/%s (文件名: %s)", hash, name)
			attachmentUrls = append(attachmentUrls, url)
		}
	}

	// 检查答案和附件是否都为空
	if answerValue == "" && len(attachmentUrls) == 0 {
		global.LOG.Info("答案和附件都为空，直接判定为0分", zap.Uint("record_id", recordID), zap.Int("index", answerIndex))

		// 获取题目详情
		questionData := gjson.Get(string(questInfo.QuestData), fmt.Sprintf("questions.%d", answerIndex))
		_ = questionData.Get("score").Int() // 获取题目满分（用于日志）

		// 直接更新为0分，不通过
		answerJSON := string(userOpenQuest.Answer)
		answerPath := fmt.Sprintf("%d", answerIndex)
		newAnswer := answerJSON
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".score", 0)
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".annotation", "未提供答案或附件")
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".open_quest_review_time", time.Now().Format("2006-01-02 15:04:05"))
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".is_ai_graded", true)

		// 检查是否所有开放题都已评分
		allReviewed := true
		for _, ans := range gjson.Get(newAnswer, "@this").Array() {
			if ans.Get("type").String() == "open_quest" {
				if !ans.Get("score").Exists() && !ans.Get("correct").Exists() {
					allReviewed = false
					break
				}
			}
		}

		// 计算总分
		var openQuestReviewStatus uint8 = 1
		var openQuestReviewTime time.Time
		var pass bool
		var userScore int64
		var openQuestScore int64

		if allReviewed {
			openQuestReviewStatus = 2
			openQuestReviewTime = time.Now()
			result, err := backend.AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(newAnswer), questInfo)
			if err == nil {
				userScore = result.UserScore
				openQuestScore = result.UserReturnScore
				pass = result.Pass
			}
		}

		// 更新数据库
		updateData := map[string]interface{}{
			"answer":                    datatypes.JSON(newAnswer),
			"open_quest_review_status":  openQuestReviewStatus,
		}
		if allReviewed {
			updateData["open_quest_review_time"] = openQuestReviewTime
			updateData["pass"] = pass
			updateData["user_score"] = userScore
			updateData["open_quest_score"] = openQuestScore
		}

		if err := global.DB.Model(&model.UserOpenQuest{}).
			Where("id = ?", recordID).
			Updates(updateData).Error; err != nil {
			global.LOG.Error("更新记录失败", zap.Error(err))
		}
		return
	}

	// 6. 获取题目详情（从quest_data中获取该索引位置的题目信息）
	questionData := gjson.Get(string(questInfo.QuestData), fmt.Sprintf("questions.%d", answerIndex))
	questionTitle := questionData.Get("title").String()
	questionScore := questionData.Get("score").Int()

	// 获取及格分数
	passingScore := gjson.Get(string(questInfo.QuestData), "passingScore").Int()

	// 7. 调用AI判题
	aiConfig, err := parseAIConfig(config.Config)
	if err != nil {
		global.LOG.Error("解析AI配置失败", zap.Error(err))
		return
	}

	systemPrompt := fmt.Sprintf(
		"你是一名web3开发专家，现在有一些web3开发学习过程中的题目和答案，需要你来进行批改，本题总分%d分，及格分%d分，请根据用户的答案和题目的要求来进行批改。如果通过，那么直接返回用户得分多少即可，如果不通过，那么需要返回不通过的分数和不通过的理由作为批注。如果答案中存在链接，那么需要访问链接中的具体网页来获取信息进行判题。如果存在附件，需要查看附件中的具体内容来进行判题。",
		questionScore,
		passingScore,
	)

	// 构建用户提示词
	userPrompt := fmt.Sprintf("题目：%s\n\n", questionTitle)

	// 如果有答案，添加答案内容
	if answerValue != "" {
		userPrompt += fmt.Sprintf("用户答案：%s", answerValue)
	} else {
		userPrompt += "用户答案：（未填写文字答案）"
	}

	// 如果有附件，添加附件信息
	if len(attachmentUrls) > 0 {
		userPrompt += "\n\n用户提交的附件："
		for i, url := range attachmentUrls {
			userPrompt += fmt.Sprintf("\n%d. %s", i+1, url)
		}
	}

	userPrompt += "\n\n请根据以上信息进行判题，返回格式：得分：X分\n批注：（如果不通过则说明理由）"

	aiResult, err := callAIAPI(aiConfig, systemPrompt, userPrompt)
	if err != nil {
		global.LOG.Error("调用AI API失败", zap.Error(err))
		return
	}

	// 8. 解析AI结果
	score, annotation := parseAIResult(aiResult, int(questionScore))
	isPassed := int64(score) >= passingScore

	// 如果通过，清空批注
	finalAnnotation := ""
	if !isPassed {
		finalAnnotation = annotation
	}

	// 9. 更新答案JSON（添加AI判题结果到指定索引位置）
	newAnswer := answerJSON
	newAnswer, _ = sjson.Set(newAnswer, answerPath+".score", score)
	newAnswer, _ = sjson.Set(newAnswer, answerPath+".annotation", finalAnnotation)
	newAnswer, _ = sjson.Set(newAnswer, answerPath+".open_quest_review_time", time.Now().Format("2006-01-02 15:04:05"))
	newAnswer, _ = sjson.Set(newAnswer, answerPath+".is_ai_graded", true)

	// 10. 检查是否所有开放题都已评分
	allReviewed := true
	for _, ans := range gjson.Get(newAnswer, "@this").Array() {
		if ans.Get("type").String() == "open_quest" {
			if !ans.Get("score").Exists() && !ans.Get("correct").Exists() {
				allReviewed = false
				break
			}
		}
	}

	// 11. 计算总分（如果所有题都已评分）
	var openQuestReviewStatus uint8 = 1 // 默认未审核
	var openQuestReviewTime time.Time
	var pass bool
	var userScore int64
	var openQuestScore int64

	if allReviewed {
		openQuestReviewStatus = 2 // 已审核
		openQuestReviewTime = time.Now()

		// 计算总分
		result, err := backend.AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(newAnswer), questInfo)
		if err == nil {
			userScore = result.UserScore
			openQuestScore = result.UserReturnScore
			pass = result.Pass
		}
	}

	// 12. 更新数据库
	updates := map[string]interface{}{
		"answer": newAnswer,
	}

	if allReviewed {
		updates["open_quest_review_status"] = openQuestReviewStatus
		updates["open_quest_review_time"] = openQuestReviewTime
		updates["open_quest_score"] = openQuestScore
		updates["pass"] = pass
		updates["user_score"] = userScore
	}

	if err := global.DB.Model(&model.UserOpenQuest{}).Where("id = ?", recordID).Updates(updates).Error; err != nil {
		global.LOG.Error("更新记录失败", zap.Error(err), zap.Uint("record_id", recordID))
		return
	}

	global.LOG.Info("AI判题完成",
		zap.Uint("record_id", recordID),
		zap.Int("answer_index", answerIndex),
		zap.Int("score", score),
		zap.Bool("all_reviewed", allReviewed),
	)
}

// calculateTotalScore 计算用户总分
func calculateTotalScore(tokenId string, openQuestScore int, questInfo model.Quest) int64 {
	// 这里需要根据题目的其他分数计算总分
	// 简化处理：假设只有开放题分数
	return int64(openQuestScore)
}

// GetAiGradeHistory 获取AI判题历史
func GetAiGradeHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var total int64

	// 查询所有已经AI判题的记录（答案中有score且有open_quest_review_time且is_ai_graded为true）
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
			user_open_quest.deleted_at IS NULL
			AND quest.status = 1
			AND json_element->>'type' = 'open_quest'
			AND json_element->>'score' IS NOT NULL
			AND json_element->>'open_quest_review_time' IS NOT NULL
			AND json_element->>'is_ai_graded' = 'true'
	`

	if err := global.DB.Raw(countSQL).Scan(&total).Error; err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}

	type HistoryItem struct {
		RecordID         uint      `json:"record_id"`
		TokenId          string    `json:"token_id"`
		Address          string    `json:"address"`
		AnswerIndex      int       `json:"answer_index"`
		QuestionTitle    string    `json:"question_title"`
		ChallengeTitle   string    `json:"challenge_title"`
		UserAnswer       string    `json:"user_answer"`
		Score            int       `json:"score"`
		Annotation       string    `json:"annotation"`
		ReviewTime       string    `json:"review_time"` // 改为字符串，因为格式不统一
		QuestionScore    int       `json:"question_score"`
		Pass             bool      `json:"pass"`
		AttachmentUrls   string    `json:"attachment_urls"` // JSON字符串
	}

	var historyItems []HistoryItem

	offset := (page - 1) * pageSize
	dataSQL := `
		SELECT
			user_open_quest.id AS record_id,
			user_open_quest.token_id,
			user_open_quest.address,
			(t.idx::int - 1) AS answer_index,
			(quest.quest_data->'questions')->(t.idx::int - 1)->>'title' AS question_title,
			quest.title AS challenge_title,
			json_element->>'value' AS user_answer,
			COALESCE(ROUND((json_element->>'score')::numeric), 0)::int AS score,
			COALESCE(json_element->>'annotation', '') AS annotation,
			json_element->>'open_quest_review_time' AS review_time,
			COALESCE(ROUND(((quest.quest_data->'questions')->(t.idx::int - 1)->>'score')::numeric), 0)::int AS question_score,
			user_open_quest.pass,
			COALESCE(json_element->>'annex', '[]') AS attachment_urls
		FROM
			user_open_quest
		JOIN
			jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
		JOIN
			quest ON quest.token_id = user_open_quest.token_id
		WHERE
			user_open_quest.deleted_at IS NULL
			AND quest.status = 1
			AND json_element->>'type' = 'open_quest'
			AND json_element->>'score' IS NOT NULL
			AND json_element->>'open_quest_review_time' IS NOT NULL
			AND json_element->>'is_ai_graded' = 'true'
		ORDER BY user_open_quest.updated_at DESC
		LIMIT ? OFFSET ?
	`

	if err := global.DB.Raw(dataSQL, pageSize, offset).Scan(&historyItems).Error; err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(gin.H{
		"list":  historyItems,
		"total": total,
		"page":  page,
		"pageSize": pageSize,
	}, c)
}

// BatchGradePreview 批量AI判题预览（不直接提交）
func BatchGradePreview(c *gin.Context) {
	var req request.BatchGradePreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	maxCount := req.MaxCount
	if maxCount <= 0 || maxCount > 50 {
		maxCount = 10
	}

	// 获取启用的AI配置
	var config model.AiJudgeConfig
	if err := global.DB.Where("enabled = ?", true).First(&config).Error; err != nil {
		response.FailWithMessage("当前没有启用的AI配置，请先在AI配置中启用一个配置", c)
		return
	}

	// 获取待判题列表
	type PendingQuest struct {
		ID      uint   `json:"id"`
		TokenId string `json:"token_id"`
		Index   int    `json:"index"`
	}
	var pendingQuests []PendingQuest

	querySQL := `
		SELECT
			user_open_quest.id,
			user_open_quest.token_id,
			(t.idx::int - 1) AS index
		FROM
			user_open_quest
		JOIN
			jsonb_array_elements(user_open_quest.answer) WITH ORDINALITY AS t(json_element, idx) ON true
		JOIN
			quest ON quest.token_id = user_open_quest.token_id
		WHERE
			user_open_quest.deleted_at IS NULL
			AND quest.status = 1
			AND json_element->>'type' = 'open_quest'
			AND json_element->>'score' IS NULL
			AND json_element->>'correct' IS NULL
		ORDER BY user_open_quest.id ASC
		LIMIT ?
	`

	if err := global.DB.Raw(querySQL, maxCount).Scan(&pendingQuests).Error; err != nil {
		response.FailWithMessage("获取待判题列表失败", c)
		return
	}

	if len(pendingQuests) == 0 {
		response.OkWithMessage("没有待判题的开放题", c)
		return
	}

	// 同步处理并返回结果（不提交到数据库）
	type PreviewResult struct {
		RecordID       uint     `json:"record_id"`
		TokenId        string   `json:"token_id"`
		Address        string   `json:"address"`
		AnswerIndex    int      `json:"answer_index"`
		QuestionTitle  string   `json:"question_title"`
		ChallengeTitle string   `json:"challenge_title"`
		UserAnswer     string   `json:"user_answer"`
		Score          int      `json:"score"`
		Annotation     string   `json:"annotation"`
		QuestionScore  int      `json:"question_score"`
		PassScore      int64    `json:"pass_score"`
		SystemPrompt   string   `json:"system_prompt"`
		UserPrompt     string   `json:"user_prompt"`
		RawResult      string   `json:"raw_result"`
		AttachmentUrls []string `json:"attachment_urls"`
	}

	var results []PreviewResult

	// 解析AI配置
	aiConfig, err := parseAIConfig(config.Config)
	if err != nil {
		response.FailWithMessage("AI配置格式错误", c)
		return
	}

	for _, pq := range pendingQuests {
		// 获取记录
		var userOpenQuest model.UserOpenQuest
		if err := global.DB.Where("id = ?", pq.ID).First(&userOpenQuest).Error; err != nil {
			continue
		}

		// 获取题目信息
		var questInfo model.Quest
		if err := global.DB.Where("token_id = ?", pq.TokenId).First(&questInfo).Error; err != nil {
			continue
		}

		// 提取答案内容
		answerJSON := string(userOpenQuest.Answer)
		answerPath := fmt.Sprintf("%d", pq.Index)
		answerData := gjson.Get(answerJSON, answerPath)
		answerValue := answerData.Get("value").String()

		// 提取附件URL
		var attachmentUrls []string
		annexArray := answerData.Get("annex").Array()
		for _, annex := range annexArray {
			hash := annex.Get("hash").String()
			name := annex.Get("name").String()
			if hash != "" {
				url := fmt.Sprintf("https://ipfs.decert.me/ipfs/%s (文件名: %s)", hash, name)
				attachmentUrls = append(attachmentUrls, url)
			}
		}

		// 如果答案和附件都为空，跳过
		if answerValue == "" && len(attachmentUrls) == 0 {
			continue
		}

		// 获取题目详情
		questionData := gjson.Get(string(questInfo.QuestData), fmt.Sprintf("questions.%d", pq.Index))
		questionTitle := questionData.Get("title").String()
		questionScore := questionData.Get("score").Int()
		passingScore := gjson.Get(string(questInfo.QuestData), "passingScore").Int()

		// 构建提示词
		systemPrompt := fmt.Sprintf(
			"你是一名web3开发专家，现在有一些web3开发学习过程中的题目和答案，需要你来进行批改，本题总分%d分，及格分%d分，请根据用户的答案和题目的要求来进行批改。如果通过，那么直接返回用户得分多少即可，如果不通过，那么需要返回不通过的分数和不通过的理由作为批注。如果答案中存在链接，那么需要访问链接中的具体网页来获取信息进行判题。如果存在附件，需要查看附件中的具体内容来进行判题。",
			questionScore,
			passingScore,
		)

		// 构建用户提示词
		userPrompt := fmt.Sprintf("题目：%s\n\n", questionTitle)

		// 如果有答案，添加答案内容
		if answerValue != "" {
			userPrompt += fmt.Sprintf("用户答案：%s", answerValue)
		} else {
			userPrompt += "用户答案：（未填写文字答案）"
		}

		// 如果有附件，添加附件信息
		if len(attachmentUrls) > 0 {
			userPrompt += "\n\n用户提交的附件："
			for i, url := range attachmentUrls {
				userPrompt += fmt.Sprintf("\n%d. %s", i+1, url)
			}
		}

		userPrompt += "\n\n请根据以上信息进行判题，返回格式：得分：X分\n批注：（如果不通过则说明理由）"

		// 调用AI API
		aiResult, err := callAIAPI(aiConfig, systemPrompt, userPrompt)
		if err != nil {
			global.LOG.Error("调用AI API失败", zap.Error(err))
			continue
		}

		// 解析AI结果
		score, annotation := parseAIResult(aiResult, int(questionScore))
		isPassed := int64(score) >= passingScore

		// 如果通过，清空批注
		finalAnnotation := ""
		if !isPassed {
			finalAnnotation = annotation
		}

		results = append(results, PreviewResult{
			RecordID:       pq.ID,
			TokenId:        pq.TokenId,
			Address:        userOpenQuest.Address,
			AnswerIndex:    pq.Index,
			QuestionTitle:  questionTitle,
			ChallengeTitle: questInfo.Title,
			UserAnswer:     answerValue,
			Score:          score,
			Annotation:     finalAnnotation,
			QuestionScore:  int(questionScore),
			PassScore:      passingScore,
			SystemPrompt:   systemPrompt,
			UserPrompt:     userPrompt,
			RawResult:      aiResult,
			AttachmentUrls: attachmentUrls,
		})

		// 延迟避免频繁调用
		time.Sleep(100 * time.Millisecond)
	}

	response.OkWithData(gin.H{
		"results": results,
		"count":   len(results),
	}, c)
}

// SubmitBatchGrade 提交批量判题结果
func SubmitBatchGrade(c *gin.Context) {
	var req request.SubmitBatchGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	if len(req.Results) == 0 {
		response.FailWithMessage("没有需要提交的判题结果", c)
		return
	}

	successCount := 0
	failCount := 0

	for _, result := range req.Results {
		// 获取记录（加锁）
		var userOpenQuest model.UserOpenQuest
		if err := global.DB.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", result.RecordID).
			First(&userOpenQuest).Error; err != nil {
			failCount++
			continue
		}

		// 检查该答案是否已经评分
		answerJSON := string(userOpenQuest.Answer)
		answerPath := fmt.Sprintf("%d", result.AnswerIndex)

		if gjson.Get(answerJSON, answerPath+".score").Exists() ||
		   gjson.Get(answerJSON, answerPath+".correct").Exists() {
			// 已评分，跳过
			failCount++
			continue
		}

		// 获取题目信息
		var questInfo model.Quest
		if err := global.DB.Where("token_id = ?", userOpenQuest.TokenId).First(&questInfo).Error; err != nil {
			failCount++
			continue
		}

		// 更新答案JSON
		newAnswer := answerJSON
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".score", result.Score)
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".annotation", result.Annotation)
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".open_quest_review_time", time.Now().Format("2006-01-02 15:04:05"))
		newAnswer, _ = sjson.Set(newAnswer, answerPath+".is_ai_graded", true)

		// 检查是否所有开放题都已评分
		allReviewed := true
		for _, ans := range gjson.Get(newAnswer, "@this").Array() {
			if ans.Get("type").String() == "open_quest" {
				if !ans.Get("score").Exists() && !ans.Get("correct").Exists() {
					allReviewed = false
					break
				}
			}
		}

		// 计算总分（如果所有题都已评分）
		var openQuestReviewStatus uint8 = 1
		var openQuestReviewTime time.Time
		var pass bool
		var userScore int64
		var openQuestScore int64

		if allReviewed {
			openQuestReviewStatus = 2
			openQuestReviewTime = time.Now()

			checkResult, err := backend.AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(newAnswer), questInfo)
			if err == nil {
				userScore = checkResult.UserScore
				openQuestScore = checkResult.UserReturnScore
				pass = checkResult.Pass
			}
		}

		// 更新数据库
		updates := map[string]interface{}{
			"answer": newAnswer,
		}

		if allReviewed {
			updates["open_quest_review_status"] = openQuestReviewStatus
			updates["open_quest_review_time"] = openQuestReviewTime
			updates["open_quest_score"] = openQuestScore
			updates["pass"] = pass
			updates["user_score"] = userScore
		}

		if err := global.DB.Model(&model.UserOpenQuest{}).Where("id = ?", result.RecordID).Updates(updates).Error; err != nil {
			failCount++
			continue
		}

		successCount++
	}

	response.OkWithData(gin.H{
		"success_count": successCount,
		"fail_count":    failCount,
		"message":       fmt.Sprintf("成功提交 %d 条，失败 %d 条", successCount, failCount),
	}, c)
}
