package request

// User login structure
type Register struct {
	Username           string `json:"username"`
	Password           string `json:"passWord"`
	Nickname           string `json:"nickname"`
	HeaderImg          string `json:"headerImg"`
	AuthorityId        string `json:"authorityId"`
	AuthoritySourceIds []uint `json:"authoritySourceIds"`
}

// User login structure
type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

// Modify password structure
type ChangePasswordStruct struct {
	ID                string `json:"id"`                                         // 用户ID
	Password          string `json:"password"`                                   // 密码
	NewPassword       string `json:"newPassword"`                                // 新密码
	RepeatNewPassword string `json:"repeatNewPassword" form:"repeatNewPassword"` // 重复新密码
}

type ResetPassword struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

type UpdateUserInfo struct {
	ID                 uint   `json:"id"`
	Nickname           string `json:"nickname"`
	HeaderImg          string `json:"headerImg"`
	AuthorityId        string `json:"authorityId"`
	AuthoritySourceIds []uint `json:"authoritySourceIds"`
	Password           string `json:"password"`
}

type UpdateSelfInfo struct {
	ID        uint   `json:"id"`
	Nickname  string `json:"nickname"`
	HeaderImg string `json:"headerImg"`
}
