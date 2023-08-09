package system

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/utils"
	"errors"
	_ "fmt"
	"strings"

	"gorm.io/gorm"
)

// CreateAuthority
// @description: 创建一个角色
// @param: auth model.Authority
// @return: authority model.Authority, err error
func CreateAuthority(auth model.Authority) (authority model.Authority, err error) {
	var authorityBox model.Authority
	if !errors.Is(global.DB.Where("authority_id = ?", auth.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return auth, errors.New("存在相同角色id")
	}
	err = global.DB.Create(&auth).Error
	return auth, err
}

// UpdateAuthority
// @description: 更改一个角色
// @param: auth model.SysAuthority
// @return: authority model.SysAuthority, err error
func UpdateAuthority(auth model.Authority) (authority model.Authority, err error) {
	err = global.DB.Where("authority_id = ?", auth.AuthorityId).First(&model.Authority{}).Updates(&auth).Error
	return auth, err
}

// @function: DeleteAuthority
// @description: 删除角色
// @param: auth *model.SysAuthority
// @return: err error
func DeleteAuthority(auth *model.Authority) (err error) {
	if errors.Is(global.DB.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色不存在")
	}
	if len(auth.Users) != 0 {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.DB.Where("authority_id = ?", auth.AuthorityId).First(&model.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.DB.Where("parent_id = ?", auth.AuthorityId).First(&model.Authority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}
	db := global.DB.Preload("SysBaseMenus").Where("authority_id = ?", auth.AuthorityId).First(auth)
	err = db.Unscoped().Delete(auth).Error
	if err != nil {
		return
	}
	//err = global.DB.Delete(&[]model.UseAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error
	//CasbinServiceApp.ClearCasbin(0, auth.AuthorityId)
	return err
}

func GetAuthorityList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&model.Authority{})
	err = db.Count(&total).Error
	var authority []model.Authority
	err = db.Limit(limit).Offset(offset).Find(&authority).Error
	return authority, total, err
}

// @function: GetAuthorityInfo
// @description: 获取所有角色信息
// @param: auth model.SysAuthority
// @return: sa model.SysAuthority, err error
func GetAuthorityInfo(auth model.Authority) (sa model.Authority, err error) {
	err = global.DB.Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(&sa).Error
	return sa, err
}

// @function: GetAuthoritySourceList
// @description: 分页获取数据
// @param: info request.PageInfo
// @return: list interface{}, err error
func GetAuthoritySourceList() (list interface{}, err error) {
	var res []model.AuthoritySource
	db := global.DB.Model(&model.AuthoritySource{})
	err = db.Find(&res).Error
	return res, err
}

// @function: SetDataAuthority
// @description: 设置角色资源权限
// @param: auth model.SysAuthority
// @return: error
func SetDataAuthority(auth request.SetDataAuthorityRequest) error {
	if auth.AuthorityId == "999" {
		return errors.New("不能修改超级管理员权限")
	}
	// 开始事务
	tx := global.DB.Begin()
	// 先删除后添加
	if err := tx.Model(&model.AuthorityRelate{}).Where("authority_id = ?", auth.AuthorityId).Delete(&model.AuthorityRelate{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除casbin表
	if err := tx.Exec("DELETE FROM \"casbin_rule\" WHERE v0 = ?", auth.AuthorityId).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加权限
	var relate []model.AuthorityRelate
	for _, v := range auth.AuthoritySourceId {
		var source model.AuthoritySource
		if err := tx.Model(&model.AuthoritySource{}).Where("id = ?", v).First(&source).Error; err != nil {
			tx.Rollback()
			return errors.New("资源不存在")
		}
		sourceList := strings.Split(source.ModelUrl, "&&")
		for _, vv := range sourceList {
			// 添加到Casbin
			if err := tx.Exec("INSERT INTO \"casbin_rule\" (\"ptype\", \"v0\", \"v1\", \"v2\", \"v3\", \"v4\", \"v5\") VALUES ('p', ?, ?, 'POST','','','')", auth.AuthorityId, vv).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Exec("INSERT INTO \"casbin_rule\" (\"ptype\", \"v0\", \"v1\", \"v2\", \"v3\", \"v4\", \"v5\") VALUES ('p', ?, ?, 'GET','','','')", auth.AuthorityId, vv).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		relate = append(relate, model.AuthorityRelate{AuthorityId: auth.AuthorityId, AuthoritySourceID: v})
	}
	if err := tx.Model(&model.AuthorityRelate{}).Create(&relate).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 更新Casbin规则
	global.Enforcer = Casbin()

	return nil
}

func SetUserAuthority(authorityId string, authoritySourceIds []uint) error {
	// 开始事务
	tx := global.DB.Begin()

	// 先删除后添加
	if err := tx.Model(&model.AuthorityRelate{}).Where("authority_id = ?", authorityId).Delete(&model.AuthorityRelate{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除casbin表
	if err := tx.Exec("DELETE FROM \"casbin_rule\" WHERE v0 = ?", authorityId).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加权限
	var relate []model.AuthorityRelate
	for _, v := range authoritySourceIds {
		var source model.AuthoritySource
		if err := tx.Model(&model.AuthoritySource{}).Where("id = ?", v).First(&source).Error; err != nil {
			tx.Rollback()
			return errors.New("资源不存在")
		}
		sourceList := strings.Split(source.ModelUrl, "&&")
		for _, vv := range sourceList {
			// 添加到Casbin
			if err := tx.Exec("INSERT INTO \"casbin_rule\" (\"ptype\", \"v0\", \"v1\", \"v2\", \"v3\", \"v4\", \"v5\") VALUES ('p', ?, ?, 'POST','','','')", authorityId, vv).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Exec("INSERT INTO \"casbin_rule\" (\"ptype\", \"v0\", \"v1\", \"v2\", \"v3\", \"v4\", \"v5\") VALUES ('p', ?, ?, 'GET','','','')", authorityId, vv).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		relate = append(relate, model.AuthorityRelate{AuthorityId: authorityId, AuthoritySourceID: v})
	}
	if err := tx.Model(&model.AuthorityRelate{}).Create(&relate).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 更新Casbin规则
	global.Enforcer = Casbin()

	return nil
}

// @function: GetAuthority
// @description: 获取角色权限
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error
func GetAuthority(auth request.SetDataAuthorityRequest) (list interface{}, err error) {
	var authorityResponse response.AuthorityResponse
	var res []model.AuthorityRelate
	db := global.DB.Model(&model.AuthorityRelate{}).Where("authority_id = ?", auth.AuthorityId)
	err = db.Find(&res).Error

	for _, SourceId := range res {
		authorityResponse.AuthoritySourceId = append(authorityResponse.AuthoritySourceId, SourceId.AuthoritySourceID)
	}
	authorityResponse.AuthorityId = auth.AuthorityId

	return authorityResponse, err
}

func UpdateUserInfo(q request.UpdateUserInfo) error {
	var user model.User
	var authority model.Authority

	// check if user exists
	err := global.DB.Where("id = ?", q.ID).First(&user).Error
	if err != nil {
		return errors.New("用户不存在")
	}

	// check if root admin
	if user.AuthorityId == "999" {
		return errors.New("不能修改超级管理员")
	}

	// check if authority exists
	result := global.DB.Where("authority_id = ?", q.AuthorityId).First(&authority)
	if result.RowsAffected == 0 {
		return errors.New("角色不存在")
	}

	user.Nickname = q.Nickname
	user.HeaderImg = q.HeaderImg
	user.AuthorityId = q.AuthorityId
	user.Password = utils.BcryptHash(q.Password)

	// TODO: 事务
	if len(q.AuthoritySourceIds) != 0 {
		err = SetUserAuthority(user.UUID.String(), q.AuthoritySourceIds)
		if err != nil {
			return errors.New("设置失败")
		}
	}

	return global.DB.Model(&model.User{}).Where("id = ?", q.ID).Updates(&user).Error
}

func UpdateSelfInfo(req request.UpdateSelfInfo) error {
	var user model.User

	user.Nickname = req.Nickname
	user.HeaderImg = req.HeaderImg

	return global.DB.Model(&model.User{}).Where("id = ?", req.ID).Updates(&user).Error
}
