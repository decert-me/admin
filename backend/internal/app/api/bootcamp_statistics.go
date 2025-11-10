package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetBootcampChallengeStatistics 获取训练营挑战统计
func GetBootcampChallengeStatistics(c *gin.Context) {
	var req struct {
		TagID      uint     `json:"tag_id" binding:"required"`
		Challenges []string `json:"challenges" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	type UserChallengeData struct {
		UserID  uint   `json:"user_id"`
		Address string `json:"address"`
		Name    string `json:"name"`
		Tags    string `json:"tags"`
		Title   string `json:"title"`
		Status  int    `json:"status"` // 0=未提交, 1=未通过, 2=通过
		UUID    string `json:"uuid"`
	}

	// 1. 获取指定标签下的所有用户
	type UserInfo struct {
		UserID  uint
		Address string
		Name    string
		Tags    string
	}
	var users []UserInfo
	userSQL := `
		SELECT
			u.id as user_id,
			u.address,
			u.name,
			string_agg(DISTINCT t.name, ',') as tags
		FROM users u
		INNER JOIN users_tag ut ON u.id = ut.user_id
		INNER JOIN tag t ON ut.tag_id = t.id
		WHERE ut.user_id IN (
			SELECT user_id FROM users_tag WHERE tag_id = ?
		)
		GROUP BY u.id, u.address, u.name
		ORDER BY u.id
	`
	if err := global.DB.Raw(userSQL, req.TagID).Scan(&users).Error; err != nil {
		global.LOG.Error("获取用户列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	// 2. 获取所有挑战的信息
	type QuestInfo struct {
		ID      uint
		TokenID string
		Title   string
		UUID    string
	}
	var quests []QuestInfo
	if err := global.DB.Table("quest").
		Select("id, token_id, title, uuid").
		Where("title IN ?", req.Challenges).
		Where("disabled = ?", false).
		Scan(&quests).Error; err != nil {
		global.LOG.Error("获取挑战列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	// 3. 查询用户的挑战记录
	var results []UserChallengeData

	for _, user := range users {
		for _, quest := range quests {
			// 查询该用户该挑战的记录（与挑战详情统计逻辑一致）
			type LogRecord struct {
				TokenID string
				Pass    bool
			}
			var logRecord LogRecord
			var hasRecord bool

			// 先查 user_challenge_log 表，检查是否有 pass=true 的记录
			var passCount int64
			err := global.DB.Table("user_challenge_log").
				Where("address = ? AND token_id = ? AND pass = true", user.Address, quest.TokenID).
				Count(&passCount).Error

			if err == nil && passCount > 0 {
				logRecord.Pass = true
				hasRecord = true
			}

			// 如果没有 pass=true 的记录，检查是否有任何记录
			if !hasRecord {
				err = global.DB.Table("user_challenge_log").
					Select("token_id, pass").
					Where("address = ? AND token_id = ?", user.Address, quest.TokenID).
					Order("pass DESC, id DESC").
					Limit(1).
					Scan(&logRecord).Error

				if err == nil && logRecord.TokenID != "" {
					hasRecord = true
				}
			}

			// 检查 user_open_quest 表（开放题），是否有 pass=true 的记录
			if !logRecord.Pass {
				var openQuestPassCount int64
				err = global.DB.Table("user_open_quest").
					Where("address = ? AND token_id = ? AND pass = true", user.Address, quest.TokenID).
					Count(&openQuestPassCount).Error

				if err == nil && openQuestPassCount > 0 {
					logRecord.Pass = true
					hasRecord = true
				}
			}

			// 检查是否已领取（在 user_challenges 或 zcloak_card 表中）
			var hasClaimed bool
			if hasRecord {
				var count int64
				// 检查 user_challenges 表
				global.DB.Table("user_challenges").
					Where("address = ? AND token_id = ?", user.Address, quest.TokenID).
					Count(&count)
				if count > 0 {
					hasClaimed = true
				} else {
					// 检查 zcloak_card 表
					global.DB.Table("zcloak_card").
						Where("address = ? AND quest_id = ?", user.Address, quest.ID).
						Count(&count)
					if count > 0 {
						hasClaimed = true
					}
				}
			}

			// 判断状态（与挑战详情统计的"挑战结果"逻辑一致）
			// 0=未提交（无记录）, 1=未通过（有记录但失败）, 2=通过（成功）
			var status int
			if !hasRecord {
				// 没有记录 = 未提交
				status = 0
			} else if hasClaimed || logRecord.Pass {
				// 已领取 或 pass=true = 通过（对应"成功"）
				status = 2
			} else {
				// 有记录但未通过 = 未完成（对应"失败"）
				status = 1
			}

			results = append(results, UserChallengeData{
				UserID:  user.UserID,
				Address: user.Address,
				Name:    user.Name,
				Tags:    user.Tags,
				Title:   quest.Title,
				Status:  status,
				UUID:    quest.UUID,
			})
		}
	}

	response.OkWithDetailed(results, "获取成功", c)
}
