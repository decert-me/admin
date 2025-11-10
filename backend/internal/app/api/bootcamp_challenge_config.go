package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetBootcampChallengeConfig 获取训练营挑战配置列表
func GetBootcampChallengeConfig(c *gin.Context) {
	type ConfigItem struct {
		QuestID      int    `json:"quest_id"`
		Title        string `json:"title"`
		UUID         string `json:"uuid"`
		Enabled      bool   `json:"enabled"`
		DisplayOrder int    `json:"display_order"`
	}

	var results []ConfigItem

	// 查询所有挑战及其配置状态
	sql := `
		SELECT
			q.id as quest_id,
			q.title,
			q.uuid,
			COALESCE(bc.enabled, false) as enabled,
			COALESCE(bc.display_order, 0) as display_order
		FROM quest q
		LEFT JOIN bootcamp_challenges bc ON q.id = bc.quest_id
		WHERE q.disabled = false
		ORDER BY
			COALESCE(bc.enabled, false) DESC,
			COALESCE(bc.display_order, 0) ASC,
			q.id ASC
	`

	if err := global.DB.Raw(sql).Scan(&results).Error; err != nil {
		global.LOG.Error("获取训练营挑战配置失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(results, "获取成功", c)
}

// UpdateBootcampChallengeConfig 更新训练营挑战配置
func UpdateBootcampChallengeConfig(c *gin.Context) {
	var req struct {
		QuestID      int  `json:"quest_id" binding:"required"`
		Enabled      bool `json:"enabled"`
		DisplayOrder int  `json:"display_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(response.TranslateValidationErrors(err), c)
		return
	}

	// 如果 enabled 为 false，强制 display_order 为 0
	displayOrder := req.DisplayOrder
	if !req.Enabled {
		displayOrder = 0
	}

	// 使用 UPSERT 语法
	sql := `
		INSERT INTO bootcamp_challenges (quest_id, enabled, display_order, updated_at)
		VALUES (?, ?, ?, NOW())
		ON CONFLICT (quest_id)
		DO UPDATE SET enabled = ?, display_order = ?, updated_at = NOW()
	`

	if err := global.DB.Exec(sql, req.QuestID, req.Enabled, displayOrder, req.Enabled, displayOrder).Error; err != nil {
		global.LOG.Error("更新训练营挑战配置失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// GetEnabledBootcampChallenges 获取启用的训练营挑战列表
func GetEnabledBootcampChallenges(c *gin.Context) {
	type ChallengeItem struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	}

	var results []ChallengeItem

	sql := `
		SELECT q.title, q.uuid
		FROM quest q
		INNER JOIN bootcamp_challenges bc ON q.id = bc.quest_id
		WHERE bc.enabled = true AND q.disabled = false
		ORDER BY bc.display_order ASC, q.id ASC
	`

	if err := global.DB.Raw(sql).Scan(&results).Error; err != nil {
		global.LOG.Error("获取启用的训练营挑战失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(results, "获取成功", c)
}
