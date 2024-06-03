package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"gorm.io/gorm"
)

// GetTagList 获取标签列表
func GetTagList(r request.GetTagListReq) (label []response.GetTagListRes, total int64, err error) {
	db := global.DB.Model(&model.Tag{})
	// 搜索
	if r.SearchVal != "" {
		db = db.Where("name like ?", "%"+r.SearchVal+"%")
	}
	// 计算总记录数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	// 使用子查询统计用户数量
	db = db.Select("tag.*, COUNT(users_tag.id) as user_num").
		Joins("LEFT JOIN users_tag ON users_tag.tag_id = tag.id").
		Group("tag.id")
	err = db.Order("id desc").Scopes(Paginate(r.Page, r.PageSize)).Find(&label).Error
	return
}

// GetTagInfo 获取标签详情
func GetTagInfo(r request.GetTagInfoReq) (label model.Tag, err error) {
	db := global.DB
	// 检查ID是否存在
	var tagExist model.Tag
	if err = db.Where("id = ?", r.TagID).First(&tagExist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果标签不存在，返回一个自定义错误消息
			return label, errors.New("标签不存在")
		}
		// 如果有其他错误则返回该错误
		return label, err
	}
	return tagExist, nil
}

// TagAdd 添加标签
func TagAdd(data model.Tag) error {
	db := global.DB.Model(&model.Tag{})
	if err := db.Create(&data).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return errors.New("标签已存在")
		}
		return err
	}
	return nil
}

// TagUpdate 修改标签
func TagUpdate(data model.Tag) error {
	db := global.DB.Model(&model.Tag{})
	// 检查ID是否存在
	var tagExist model.Tag
	if err := db.Where("id = ?", data.ID).First(&tagExist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果标签不存在，返回一个自定义错误消息
			return errors.New("标签不存在")
		}
		// 如果有其他错误则返回该错误
		return err
	}
	// 更新标签
	if err := db.Where("id = ?", data.ID).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

// GetTagUserList 查询标签用户列表
func GetTagUserList(r request.GetTagUserListReq) (label []response.GetTagUserListRes, total int64, err error) {
	if r.TagID == 0 {
		err = errors.New("invalid tag ID")
		return
	}
	db := global.DB.Model(&model.UsersTag{}).
		Joins("JOIN users ON users_tag.user_id = users.id").
		Where("users_tag.tag_id = ?", r.TagID)
	// 搜索
	if r.SearchVal != "" {
		db = db.Where("users.address like ?", "%"+r.SearchVal+"%")
	}
	// 计算总记录数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Select("users.id, to_timestamp(users.creation_timestamp) as created_at, users.nickname, users.address")

	// 执行查询
	err = db.Scopes(Paginate(r.Page, r.PageSize)).Find(&label).Error
	return
}

// TagUserUpdate 添加用户标签
func TagUserUpdate(r request.TagUserUpdateReq) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 判断用户是否存在
	var user model.Users
	if err := tx.Where("id = ?", r.UserID).First(&user).Error; err != nil {
		tx.Rollback()
		return errors.New("用户不存在")
	}
	// 判断标签是否存在
	var tag model.Tag
	if err := tx.Where("id = ?", r.TagID).First(&tag).Error; err != nil {
		tx.Rollback()
		return errors.New("标签不存在")
	}
	// 检查用户标签是否已经存在
	var existingUserTag model.UsersTag
	if err := tx.Where("user_id = ? AND tag_id = ?", r.UserID, r.TagID).First(&existingUserTag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果标签不存在，则创建新的标签关联
			userTag := model.UsersTag{UserID: r.UserID, TagID: r.TagID}
			if err := tx.Create(&userTag).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 其他错误情况，回滚事务
			tx.Rollback()
			return err
		}
	}

	// Update user's nickname if provided
	if r.NickName != nil {
		if err := tx.Model(&model.Users{}).Where("id = ?", r.UserID).Update("nickname", *r.NickName).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

// TagDeleteBatch 批量删除标签
func TagDeleteBatch(r request.TagDeleteBatchReq) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 删除标签
	if err := tx.Where("id in (?)", r.TagIDs).Delete(&model.Tag{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除用户标签
	if err := tx.Where("tag_id in (?)", r.TagIDs).Delete(&model.UsersTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}

// TagUserDeleteBatch 批量删除用户标签
func TagUserDeleteBatch(r request.TagUserDeleteBatchReq) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 删除用户标签
	if err := tx.Where("user_id in (?) AND tag_id = ?", r.UserIDs, r.TagID).Delete(&model.UsersTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
