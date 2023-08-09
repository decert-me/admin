package system

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/utils"
	"errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// @function: Register
// @description: 用户注册
// @param: u model.User
// @return: userInter model.User, err error
func Register(u model.User) (userInter model.User, err error) {
	var user model.User
	if !errors.Is(global.DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	// 设置用户角色权限
	if u.AuthorityId == "" {
		u.AuthorityId = "1"
	}

	// 判断角色ID是否存在
	var authority model.Authority
	if errors.Is(global.DB.Model(&model.Authority{}).Where("authority_id = ?", u.AuthorityId).First(&authority).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("角色不存在")
	}

	// 设置资源权限
	if len(u.AuthoritySourceIds) != 0 {
		// 判断是否超出权限范围
		var authorityRelates []model.AuthorityRelate
		global.DB.Where(map[string]interface{}{"authority_id": u.AuthorityId}).Find(&authorityRelates)

		isIn := true
		for _, v := range u.AuthoritySourceIds {
			innerIsIn := false
			for _, value := range authorityRelates {
				if v == value.AuthoritySourceID {
					innerIsIn = true
					continue
				}
			}

			if !innerIsIn {
				isIn = false
				break
			}
		}
		if !isIn {
			return userInter, errors.New("超出角色权限范围")
		}

		err := SetUserAuthority(u.UUID.String(), u.AuthoritySourceIds)
		if err != nil {
			return u, err
		}
	}

	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.NewV4()
	err = global.DB.Create(&u).Error

	u.Authority = authority

	return u, err
}

// Login @function: Login
// @description: 用户登录
// @param: u *model.User
// @return: userInter *model.User, err error
func Login(u *model.User) (userInter *model.User, err error) {
	var user model.User
	err = global.DB.Where("username = ?", u.Username).Preload("Authority").First(&user).Error
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("用户名或密码错误")
	}

	// 返回资源ID
	var authoritySourceIds []uint
	// 先查用户资源，再查角色资源
	_ = global.DB.Model(&model.AuthorityRelate{}).Select("authority_source_id").Where("authority_id = ?", user.Username).Find(&authoritySourceIds).Error
	if len(authoritySourceIds) == 0 {
		_ = global.DB.Model(&model.AuthorityRelate{}).Select("authority_source_id").Where("authority_id = ?", user.AuthorityId).Find(&authoritySourceIds).Error
	}
	user.AuthoritySourceIds = authoritySourceIds

	return &user, err
}

// @function: ChangePassword
// @description: 修改用户密码
// @param: u *model.User, newPassword string
// @return: userInter *model.User, err error
func ChangePassword(u *model.User, newPassword string) (userInter *model.User, err error) {
	var user model.User
	err = global.DB.Where("username = ?", u.Username).First(&user).Error
	if err != nil {
		return nil, errors.New("找不到该用户")
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.DB.Save(&user).Error
	return &user, err
}

// @function: GetUserInfoList
// @description: 分页获取所有用户数据
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error
func GetUserInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&model.User{})
	var userList []model.User
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(int(limit)).Offset(int(offset)).Preload("Authority").Order("id desc").Find(&userList).Error

	if err == nil {
		for i, v := range userList {
			// 返回资源ID
			var authoritySourceIds []uint

			// 先查用户资源，再查角色资源
			_ = global.DB.Model(&model.AuthorityRelate{}).Select("authority_source_id").Where("authority_id = ?", v.Username).Find(&authoritySourceIds).Error
			if len(authoritySourceIds) == 0 {
				_ = global.DB.Model(&model.AuthorityRelate{}).Select("authority_source_id").Where("authority_id = ?", v.AuthorityId).Find(&authoritySourceIds).Error
			}
			userList[i].AuthoritySourceIds = authoritySourceIds
		}
	}

	return userList, total, err
}

// @function: DeleteUser
// @description: 删除用户
// @param: id float64
// @return: err error
func DeleteUser(id uint) (err error) {
	var user model.User
	result := global.DB.Where("id = ?", id).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("找不到该用户")
	}

	if user.AuthorityId == "999" {
		return errors.New("不能删除超级管理员")
	}

	// TODO: 上下级限制

	err = global.DB.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	return err
}

// @function: GetUserInfo
// @description: 获取单个用户信息
// @param: ID request.IDRequest
// @return: user system.User, err error
func GetUserInfo(ID uint) (userInfo model.User, err error) {
	var user model.User
	err = global.DB.Model(&model.User{}).Where("id = ?", int(ID)).First(&user).Error
	if err != nil {
		return user, errors.New("找不到该用户")
	}

	// 返回资源ID
	var authoritySourceIds []uint
	// 先查用户资源，再查角色资源
	_ = global.DB.Model(&model.AuthorityRelate{}).Select("authority_source_id").Where("authority_id = ?", user.Username).Find(&authoritySourceIds).Error
	if len(authoritySourceIds) == 0 {
		_ = global.DB.Model(&model.AuthorityRelate{}).Select("authority_source_id").Where("authority_id = ?", user.AuthorityId).Find(&authoritySourceIds).Error
	}
	user.AuthoritySourceIds = authoritySourceIds

	return user, err
}

// @function: resetPassword
// @description: 管理员重置用户密码
// @param: ID uint, pass string
// @return: err error
func ResetPassword(operatorUID uint, q request.ResetPassword) (err error) {
	var targetUser model.User
	result := global.DB.Where("id = ?", q.ID).First(&targetUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("找不到该用户")
	}

	if targetUser.AuthorityId == "999" {
		return errors.New("不能重置超级管理员密码")
	}

	var operator model.User
	_ = global.DB.Where("id = ?", operatorUID).First(&operator)

	if targetUser.AuthorityId >= operator.AuthorityId {
		return errors.New("不能修改级同/高级别的用户信息")
	}

	err = global.DB.Model(&model.User{}).Where("id = ?", q.ID).Update("password", utils.BcryptHash(q.Password)).Error

	return err
}
