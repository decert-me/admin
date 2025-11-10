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
		UserID    uint   `json:"user_id"`
		Address   string `json:"address"`
		Name      string `json:"name"`
		Tags      string `json:"tags"`
		Title     string `json:"title"`
		UserScore int    `json:"user_score"`
		Status    int    `json:"status"`
		UUID      string `json:"uuid"`
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
				TokenID   string
				UserScore int64
				Pass      bool
			}
			var logRecord LogRecord
			var hasRecord bool

			// 先查 user_challenge_log 表，读取 pass 字段和 user_score
			err := global.DB.Table("user_challenge_log").
				Select("token_id, user_score, pass").
				Where("address = ? AND token_id = ?", user.Address, quest.TokenID).
				Order("pass DESC, user_score DESC, id DESC").
				Limit(1).
				Scan(&logRecord).Error

			if err == nil && logRecord.TokenID != "" {
				hasRecord = true
			} else if err != nil && err.Error() != "record not found" {
				global.LOG.Error("查询挑战记录失败!", zap.Error(err))
			}

			// 如果是开放题，可能在 user_open_quest 表中
			if !hasRecord {
				logRecord = LogRecord{} // 重置
				err = global.DB.Table("user_open_quest").
					Select("token_id, user_score, pass").
					Where("address = ? AND token_id = ?", user.Address, quest.TokenID).
					Order("pass DESC, user_score DESC, id DESC").
					Limit(1).
					Scan(&logRecord).Error

				if err == nil && logRecord.TokenID != "" {
					hasRecord = true
				}
			}

			// 检查是否在 user_challenges 表中（已领取NFT）
			var hasClaimed bool
			if hasRecord {
				var count int64
				global.DB.Table("user_challenges").
					Where("address = ? AND token_id = ?", user.Address, quest.TokenID).
					Count(&count)
				if count > 0 {
					hasClaimed = true
				}
			}

			// 如果找不到记录，标记为未完成
			if !hasRecord {
				results = append(results, UserChallengeData{
					UserID:    user.UserID,
					Address:   user.Address,
					Name:      user.Name,
					Tags:      user.Tags,
					Title:     quest.Title,
					UserScore: 0,
					Status:    0, // 未完成
					UUID:      quest.UUID,
				})
				continue
			}

			// 转换分数（从10000分制转为100分制）
			scorePercent := int(logRecord.UserScore / 100)

			// 判断是否通过（与挑战详情统计逻辑一致）
			// 只要 pass=true 或者已领取NFT，就标记为通过
			status := 0 // 未通过
			if logRecord.Pass || hasClaimed {
				status = 2 // 通过
			}

			results = append(results, UserChallengeData{
				UserID:    user.UserID,
				Address:   user.Address,
				Name:      user.Name,
				Tags:      user.Tags,
				Title:     quest.Title,
				UserScore: scorePercent,
				Status:    status,
				UUID:      quest.UUID,
			})
		}
	}

	response.OkWithDetailed(results, "获取成功", c)
}
