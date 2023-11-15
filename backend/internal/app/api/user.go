package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/system"
	"backend/internal/app/utils"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.User) {
	j := &utils.JWT{SigningKey: []byte(global.CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(utils.BaseClaims{
		ID:          user.ID,
		Address:     user.Address,
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
	user := &model.User{Username: r.Username, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Address: r.Address}
	userReturn, err := system.Register(utils.GetUserID(c), *user)
	if err != nil {
		global.LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(response.UserResponse{User: userReturn}, err.Error(), c)
	} else {
		response.OkWithDetailed(response.UserResponse{User: userReturn}, "添加成功", c)
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

	if err := system.DeleteUser(uid, reqId.ID); err != nil {
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
	//uid := utils.GetUserID(c)
	uid := c.Query("id")
	ReqUser, err := system.GetUserInfo(cast.ToUint(uid))

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
	if err := system.UpdateUserInfo(utils.GetUserID(c), q); err != nil {
		global.LOG.Error("编辑失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("编辑成功", c)
	}
}

func UpdateSelfInfo(c *gin.Context) {
	var q request.UpdateSelfInfo
	_ = c.ShouldBindJSON(&q)

	if err := system.UpdateSelfInfo(utils.GetUserID(c), q); err != nil {
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

// GetLoginMessage
// @Tags SignApi
// @Summary 获取登录签名消息
// @accept application/json
// @Produce application/json
// @Router /sign/getLoginMessage [post]
func GetLoginMessage(c *gin.Context) {
	var request request.GetLoginMessageRequest
	_ = c.ShouldBindQuery(&request)
	if loginMessage, err := system.GetLoginMessage(request.Address); err != nil {
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(map[string]string{"loginMessage": loginMessage}, "获取成功", c)
	}
}

// AuthLoginSign
// @Tags SignApi
// @Summary 校验登录签名
// @accept application/json
// @Produce application/json
// @Router /sign/authLoginSign [post]
func AuthLoginSign(c *gin.Context) {
	var request request.AuthLoginSignRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.FailWithMessage("参数错误："+err.Error(), c)
		return
	}
	if user, err := system.AuthLoginSignRequest(request); err != nil {
		response.FailWithMessage("暂无权限，请联系管理员", c)
	} else {
		tokenNext(c, user)
	}
}
