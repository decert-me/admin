package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/system"
	"backend/internal/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"go.uber.org/zap"
)

// @Tags User
// @Summary 用户登录
// @Produce  application/json
// @Param data body model.UserLogin true "用户名, 密码, 验证码"
// @Success 200 {object} response.Response{data=request.LoginResponse,msg=string} "返回包括用户信息,token,过期时间"
// @Router /base/login [post]
func Login(c *gin.Context) {
	var l request.Login
	_ = c.ShouldBindJSON(&l)
	// 校验登录信息是否完整
	if global.CONFIG.System.UseCaptcha {
		// 开启验证码
		if err := utils.Verify(l, utils.LoginVerify); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	} else {
		// 不使用验证码
		if err := utils.Verify(l, utils.LoginVerifyNotCaptcha); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	// 登录信息
	var userLogin model.UserLogin
	userLogin.Username = l.Username
	userLogin.IP = c.ClientIP()
	userLogin.Agent = c.Request.UserAgent()
	ua := user_agent.New(c.Request.UserAgent())
	userLogin.LoginLocation = "" // 取消使用QQWRT
	userLogin.OS = ua.OS()
	userLogin.Browser, _ = ua.Browser()

	// 校验验证码
	if global.CONFIG.System.UseCaptcha {
		if store.Verify(l.CaptchaId, l.Captcha, true) {
			u := &model.User{Username: l.Username, Password: l.Password}
			if user, err := system.Login(u); err != nil {
				userLogin.Status = 0
				userLogin.ErrorMessage = err.Error()
				// 写入登录日志
				system.CreateUserLogin(userLogin)
				global.LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
				response.FailWithMessage("用户名不存在或者密码错误", c)
			} else {
				userLogin.Status = 1
				// 写入登录日志
				system.CreateUserLogin(userLogin)
				tokenNext(c, *user)
			}
		} else {
			userLogin.Status = 0
			userLogin.ErrorMessage = "验证码错误"
			system.CreateUserLogin(userLogin)
			response.FailWithMessage("验证码错误", c)
		}
	} else {
		// 不使用验证码
		u := &model.User{Username: l.Username, Password: l.Password}
		if user, err := system.Login(u); err != nil {
			userLogin.Status = 0
			userLogin.ErrorMessage = err.Error()
			// 写入登录日志
			system.CreateUserLogin(userLogin)
			global.LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
			response.FailWithMessage("用户名不存在或者密码错误", c)
		} else {
			userLogin.Status = 1
			// 写入登录日志
			system.CreateUserLogin(userLogin)
			tokenNext(c, *user)
		}
	}
}

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.User) {
	j := &utils.JWT{SigningKey: []byte(global.CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(utils.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		Nickname:    user.Nickname,
		UserName:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	response.OkWithDetailed(response.LoginResponse{
		User:  user,
		Token: token,
	}, "登录成功", c)
	return
}

// @Tags User
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body request.Register true "用户名, 昵称, 密码, 角色ID"
// @Success 200 {object} response.Response{data=response.UserResponse,msg=string} "用户注册账号,返回包括用户信息"
// @Router /user/register [post]
func Register(c *gin.Context) {
	var r request.Register
	_ = c.ShouldBindJSON(&r)
	if err := utils.Verify(r, utils.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := &model.User{Username: r.Username, Nickname: r.Nickname, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, AuthoritySourceIds: r.AuthoritySourceIds}
	userReturn, err := system.Register(*user)
	if err != nil {
		global.LOG.Error("注册失败!", zap.Error(err))
		response.FailWithDetailed(response.UserResponse{User: userReturn}, err.Error(), c)
	} else {
		response.OkWithDetailed(response.UserResponse{User: userReturn}, "注册成功", c)
	}
}

// @Tags User
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.ChangePasswordStruct true "用户名, 原密码, 新密码"
// @Success 200 {object} response.Response{msg=string} "用户修改密码"
// @Router /user/changePassword [post]
func ChangePassword(c *gin.Context) {
	var user request.ChangePasswordStruct
	_ = c.ShouldBindJSON(&user)
	if err := utils.Verify(user, utils.ChangePasswordVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	username := utils.GetUsername(c)

	if user.NewPassword != user.RepeatNewPassword {
		response.FailWithMessage("修改失败：前后密码不一致", c)
		return
	}

	u := &model.User{Username: username, Password: user.Password}
	if _, err := system.ChangePassword(u, user.NewPassword); err != nil {
		global.LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败：原密码与当前账户不符", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// @Tags User
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindQuery(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := system.GetUserInfoList(pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags User
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IDRequest true "用户ID"
// @Success 200 {object} response.Response{msg=string} "删除用户"
// @Router /user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var reqId request.IDRequest
	_ = c.ShouldBindJSON(&reqId)

	if err := utils.Verify(reqId, utils.IDVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	uid := utils.GetUserID(c)
	if uid == uint(reqId.ID) {
		response.FailWithMessage("删除失败, 自杀失败", c)
		return
	}

	if err := system.DeleteUser(reqId.ID); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags User
// @Summary 获取用户个人信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取用户信息"
// @Router /user/getUserInfo [get]
func GetSelfInfo(c *gin.Context) {
	uid := utils.GetUserID(c)

	ReqUser, err := system.GetUserInfo(uid)

	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.UserInfo{
		User: ReqUser,
	}, "获取成功", c)
}

func UpdateUserInfo(c *gin.Context) {
	var q request.UpdateUserInfo
	_ = c.ShouldBindJSON(&q)
	if err := utils.Verify(q, utils.IDVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := system.UpdateUserInfo(q); err != nil {
		global.LOG.Error("编辑失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("编辑成功", c)
	}
}

func UpdateSelfInfo(c *gin.Context) {
	var q request.UpdateSelfInfo
	_ = c.ShouldBindJSON(&q)

	q.ID = utils.GetUserID(c)

	if err := system.UpdateSelfInfo(q); err != nil {
		global.LOG.Error("编辑失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
	} else {
		response.OkWithMessage("编辑成功", c)
	}
}

// @Tags User
// @Summary 重置成员密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body system.ResetPassword true "ID"
// @Success 200 {object} response.Response{msg=string} "重置成员密码"
// @Router /user/resetPassword [post]
func ResetPassword(c *gin.Context) {
	var q request.ResetPassword
	_ = c.ShouldBindJSON(&q)
	if len(q.Password) < 6 {
		response.FailWithMessage("密码长度必须不小于6位", c)
		return
	}

	operatorID := utils.GetUserID(c)

	if err := system.ResetPassword(operatorID, q); err != nil {
		global.LOG.Error("重置失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("重置成功", c)
	}
}
