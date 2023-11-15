package request

// User login structure
type Register struct {
	Username    string `json:"username"`
	Address     string `json:"address"`
	Password    string `json:"passWord"`
	HeaderImg   string `json:"headerImg"`
	AuthorityId string `json:"authorityId"`
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
	ID uint `json:"id"`
	//Nickname    string `json:"nickname"`
	UserName    string `json:"userName"`
	Address     string `json:"address"`
	HeaderImg   string `json:"headerImg"`
	AuthorityId string `json:"authorityId"`
	//Password    string `json:"password"`
}

type UpdateSelfInfo struct {
	ID        string `json:"id"`
	UserName  string `json:"userName"`
	Address   string `json:"address"`
	HeaderImg string `json:"headerImg"`
}

type GetLoginMessageRequest struct {
	Address string `json:"address" form:"address"`
}

type AuthLoginSignRequest struct {
	Address   string `json:"address" form:"address" binding:"required"`
	Message   string `json:"message" form:"message" binding:"required"`
	Signature string `json:"signature" form:"signature" binding:"required"`
}
