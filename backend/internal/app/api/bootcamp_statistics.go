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

	var results []UserChallengeData

	// 查询SQL：获取指定标签下的所有用户及其挑战完成情况
	sql := `
		SELECT
			u.id as user_id,
			u.address,
			u.name,
			string_agg(DISTINCT t.name, ',') as tags,
			q.title,
			q.uuid,
			COALESCE(uc.user_score, 0) as user_score,
			COALESCE(uc.status, 0) as status
		FROM users u
		INNER JOIN users_tag ut ON u.id = ut.user_id
		INNER JOIN tag t ON ut.tag_id = t.id
		LEFT JOIN user_challenges uc ON u.address = uc.address
		LEFT JOIN quest q ON uc.token_id = q.token_id AND q.title = ANY($2::text[])
		WHERE ut.user_id IN (
			SELECT user_id FROM users_tag WHERE tag_id = $1
		)
		GROUP BY u.id, u.address, u.name, q.title, q.uuid, uc.user_score, uc.status
		ORDER BY u.id
	`

	// 将challenges转换为PostgreSQL数组格式
	challengesArray := "{" + `"` + req.Challenges[0] + `"`
	for i := 1; i < len(req.Challenges); i++ {
		challengesArray += `,"` + req.Challenges[i] + `"`
	}
	challengesArray += "}"

	if err := global.DB.Raw(sql, req.TagID, challengesArray).Scan(&results).Error; err != nil {
		global.LOG.Error("获取训练营挑战统计失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(results, "获取成功", c)
}
