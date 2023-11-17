package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"errors"
	"gorm.io/gorm"
)

// GetOpenQuestPermList 获取权限列表
func GetOpenQuestPermList(r request.GetOpenQuestPermListRequest) (data []*model.OpenQuestPerm, total int64, err error) {
	db := global.DB.Model(&model.OpenQuestPerm{})
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Scopes(Paginate(r.Page, r.PageSize)).Find(&data).Error
	return
}

// AddOpenQuestPerm 添加权限
func AddOpenQuestPerm(r request.AddOpenQuestPermRequest) (err error) {
	perm := model.OpenQuestPerm{
		Address: r.Address,
	}
	err = global.DB.Create(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = errors.New("地址权限已存在")
		}
	}
	return
}

// DeleteOpenQuestPerm 删除权限
func DeleteOpenQuestPerm(r request.DeleteOpenQuestPermRequest) (err error) {
	err = global.DB.Delete(&model.OpenQuestPerm{}, "address ILIKE ?", r.Address).Error
	return
}
