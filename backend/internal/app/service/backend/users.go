package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
)

// GetUsersList 获取用户列表
func GetUsersList(r request.GetUsersListReq) (label []response.GetUsersListRes, total int64, err error) {
	// 构建基础查询
	db := global.DB.Table("users").
		Select("users.id as user_id, users.address, users.name, string_agg(tag.name, ',') as tags, to_timestamp(users.creation_timestamp) as created_at").
		Joins("LEFT JOIN users_tag on users_tag.user_id = users.id").
		Joins("LEFT JOIN tag on tag.id = users_tag.tag_id").
		Group("users.id, users.address")

	// 应用搜索条件
	if r.SearchTag != "" {
		db = db.Where("tag.name LIKE ?", "%"+r.SearchTag+"%")
	}
	if r.SearchAddress != "" {
		db = db.Where("(users.address ILIKE ? OR users.name ILIKE ?)", "%"+r.SearchAddress+"%", "%"+r.SearchAddress+"%")
	}
	// 获取总数用于分页
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	// 执行查询
	err = db.Scopes(Paginate(r.Page, r.PageSize)).Find(&label).Error
	return
}

// GetUsersInfo 查询用户详情
func GetUsersInfo(r request.GetUsersInfoReq) (label response.GetUsersInfoRes, err error) {
	// 查询用户基本信息
	if err := global.DB.Table("users").
		Select("id as user_id, address, name, to_timestamp(creation_timestamp) as created_at").
		Where("id = ?", r.UserID).
		Limit(1).
		Find(&label).Error; err != nil {
		return label, err
	}

	// 查询用户所拥有的标签信息
	var tags []model.Tag
	if err := global.DB.Table("tag").
		Joins("LEFT JOIN users_tag on tag.id = users_tag.tag_id").
		Where("users_tag.user_id = ?", r.UserID).
		Find(&tags).Error; err != nil {
		return label, err
	}

	// 构造返回结果
	label.Tag = tags

	return label, nil
}

// UpdateUsersInfo 修改用户
func UpdateUsersInfo(r request.UpdateUsersInfoReq) error {
	// 开启事务
	tx := global.DB.Begin()

	// 删除该用户当前所有的标签关联
	if err := tx.Where("user_id = ?", r.UserID).Delete(&model.UsersTag{}).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 为用户添加新的标签关联
	for _, tagID := range r.TagIds {
		// 查询标签是否存在
		var tag model.Tag
		if err := tx.Where("id = ?", tagID).First(&tag).Error; err != nil {
			tx.Rollback() // 回滚事务
			return errors.New("标签不存在")
		}
		userTag := model.UsersTag{UserID: r.UserID, TagID: tagID}
		if err := tx.Create(&userTag).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}
	// 更新用户信息
	if r.Name != nil {
		if err := tx.Model(&model.Users{}).Where("id = ?", r.UserID).Update("name", r.Name).Error; err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}
	// 提交事务
	return tx.Commit().Error
}
