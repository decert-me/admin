package request

// CreateAiJudgeConfigRequest 创建AI判题配置请求
type CreateAiJudgeConfigRequest struct {
	Title  string `json:"title" binding:"required"`
	Config string `json:"config" binding:"required"`
}

// UpdateAiJudgeConfigRequest 更新AI判题配置请求
type UpdateAiJudgeConfigRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Config string `json:"config" binding:"required"`
}

// ToggleAiJudgeConfigRequest 切换AI判题配置启用状态请求
type ToggleAiJudgeConfigRequest struct {
	ID uint `json:"id" binding:"required"`
}

// DeleteAiJudgeConfigRequest 删除AI判题配置请求
type DeleteAiJudgeConfigRequest struct {
	ID uint `json:"id" binding:"required"`
}

// AiGradeRequest AI判题请求
type AiGradeRequest struct {
	QuestionTitle  string   `json:"question_title" binding:"required"`  // 题目标题
	QuestionScore  int      `json:"question_score" binding:"required"`  // 题目总分
	PassScore      int      `json:"pass_score" binding:"required"`      // 及格分
	UserAnswer     string   `json:"user_answer" binding:"required"`     // 用户答案
	AttachmentUrls []string `json:"attachment_urls"`                    // 附件URL列表
}

// BatchGradeRequest 批量AI判题请求
type BatchGradeRequest struct {
	MaxCount int `json:"max_count"` // 最多处理多少条，默认10
}

// ToggleAutoGradingRequest 切换自动判题状态请求
type ToggleAutoGradingRequest struct {
	ID          uint `json:"id" binding:"required"`
	AutoGrading bool `json:"auto_grading"`
}

// GetAiGradeHistoryRequest 获取AI判题历史请求
type GetAiGradeHistoryRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

// BatchGradePreviewRequest 批量AI判题预览请求（不直接提交）
type BatchGradePreviewRequest struct {
	MaxCount int `json:"max_count"` // 最多处理多少条，默认10
}

// SubmitBatchGradeRequest 提交批量判题结果请求
type SubmitBatchGradeRequest struct {
	Results []SubmitGradeResult `json:"results" binding:"required"`
}

// SubmitGradeResult 单个判题结果
type SubmitGradeResult struct {
	RecordID    uint   `json:"record_id" binding:"required"`    // user_open_quest记录ID
	AnswerIndex int    `json:"answer_index" binding:"required"` // 答案数组索引
	Score       int    `json:"score" binding:"required"`        // 分数
	Annotation  string `json:"annotation"`                      // 批注
}
