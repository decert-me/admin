package api

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/service/system"
	"backend/internal/app/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityApi struct{}

// CreateAuthority
// @Summary 创建角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.Authority true "权限id, 权限名, 父角色id"
// @Success 200 {object} response.Response{data=systemRes.AuthorityResponse,msg=string} "创建角色,返回包括系统角色详情"
// @Router /authority/createAuthority [post]
func CreateAuthority(c *gin.Context) {
	var authority model.Authority
	_ = c.ShouldBindJSON(&authority)
	if err := utils.Verify(authority, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if authBack, err := system.CreateAuthority(authority); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		_ = system.UpdateCasbin(authority.AuthorityId, request.DefaultCasbin())
		response.OkWithDetailed(response.SysAuthorityResponse{Authority: authBack}, "创建成功", c)
	}
}

// DeleteAuthority
// @Summary 删除角色
// @accept application/json
// @Produce application/json
// @Param data body model.Authority true "删除角色"
// @Success 200 {object} response.Response{msg=string} "删除角色"
// @Router /authority/deleteAuthority [post]
func DeleteAuthority(c *gin.Context) {
	var authority model.Authority
	_ = c.ShouldBindJSON(&authority)
	if err := utils.Verify(authority, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := system.DeleteAuthority(&authority); err != nil { // 删除角色之前需要判断是否有用户正在使用此角色
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateAuthority
// @Tags Authority
// @Summary 更新角色信息
// @Produce application/json
// @Param data body model.Authority true "权限id, 权限名, 父角色id"
// @Success 200 {object} response.Response{data=systemRes.SysAuthorityResponse,msg=string} "更新角色信息,返回包括系统角色详情"
// @Router /authority/updateAuthority [post]
func UpdateAuthority(c *gin.Context) {
	var auth model.Authority
	_ = c.ShouldBindJSON(&auth)
	if err := utils.Verify(auth, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if authority, err := system.UpdateAuthority(auth); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.SysAuthorityResponse{Authority: authority}, "更新成功", c)
	}
}

func GetAuthorityList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := system.GetAuthorityList(pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// SetDataAuthority @Tags Authority
// @Summary 设置角色资源权限
// @accept application/json
// @Produce application/json
// @Param data body model.Authority true "设置角色资源权限"
// @Success 200 {object} response.Response{msg=string} "设置角色资源权限"
// @Router /authority/setDataAuthority [post]
func SetDataAuthority(c *gin.Context) {
	var auth request.SetDataAuthorityRequest
	_ = c.ShouldBindJSON(&auth)
	if err := utils.Verify(auth, utils.AuthorityDataVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := system.SetDataAuthority(auth); err != nil {
		global.LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败，"+err.Error(), c)
	} else {
		response.OkWithMessage("设置成功", c)
	}
}

// @Tags Authority
// @Summary 分页获取角色列表

// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取角色列表,返回包括列表,总数,页码,每页数量"
// @Router /authority/getAuthorityList [post]
func GetAuthoritySourceList(c *gin.Context) {
	if list, err := system.GetAuthoritySourceList(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(list, "获取成功", c)
	}
}

func GetAuthority(c *gin.Context) {
	var auth request.SetDataAuthorityRequest
	_ = c.ShouldBindQuery(&auth)
	if err := utils.Verify(auth, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, err := system.GetAuthority(auth); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(list, "获取成功", c)
	}
}
