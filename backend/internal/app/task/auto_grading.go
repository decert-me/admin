package task

import (
	"backend/internal/app/api"
	"backend/internal/app/global"
	"backend/internal/app/model"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var AutoGradingCron *cron.Cron

// StartAutoGradingTask 启动自动判题定时任务
func StartAutoGradingTask() {
	AutoGradingCron = cron.New(cron.WithSeconds()) // 支持秒级别

	// 每分钟执行一次
	_, err := AutoGradingCron.AddFunc("0 * * * * *", func() {
		runAutoGrading()
	})

	if err != nil {
		global.LOG.Error("添加自动判题任务失败", zap.Error(err))
		return
	}

	AutoGradingCron.Start()
	global.LOG.Info("自动判题定时任务已启动")
}

// StopAutoGradingTask 停止自动判题定时任务
func StopAutoGradingTask() {
	if AutoGradingCron != nil {
		AutoGradingCron.Stop()
		global.LOG.Info("自动判题定时任务已停止")
	}
}

// runAutoGrading 执行自动判题
func runAutoGrading() {
	// 检查是否有启用自动判题的配置
	var config model.AiJudgeConfig
	if err := global.DB.Where("enabled = ? AND auto_grading = ?", true, true).First(&config).Error; err != nil {
		// 没有启用自动判题
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
		LEFT JOIN
			users ON users.address = user_open_quest.address
		LEFT JOIN
			users_tag ON users_tag.user_id = users.id
		WHERE
			user_open_quest.deleted_at IS NULL
			AND quest.status = 1
			AND json_element->>'type' = 'open_quest'
			AND json_element->>'score' IS NULL
			AND json_element->>'correct' IS NULL
			AND users_tag.user_id IS NULL
		ORDER BY user_open_quest.id ASC
		LIMIT 10
	`

	if err := global.DB.Raw(querySQL).Scan(&pendingQuests).Error; err != nil {
		global.LOG.Error("获取待判题列表失败", zap.Error(err))
		return
	}

	if len(pendingQuests) == 0 {
		return // 没有待判题的题目
	}

	global.LOG.Info("开始自动判题", zap.Int("count", len(pendingQuests)))

	// 批量判题
	for _, pq := range pendingQuests {
		api.ProcessOneOpenQuest(pq.ID, pq.TokenId, pq.Index, config)
		// 延迟100ms，避免频繁调用AI API
		time.Sleep(100 * time.Millisecond)
	}
}
