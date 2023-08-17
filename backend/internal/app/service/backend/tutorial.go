package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"errors"
)

func GetTutorialList(info request.GetTutorialListStatusRequest) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&model.Tutorial{})
	// 语言
	if info.Language != 0 {
		db.Where("language = ?", info.Language)
	}
	// 状态
	if info.Status != 0 {
		db.Where("status = ?", info.Status)
	}
	// 根据分类要求过滤
	if info.Category != nil && len(info.Category) != 0 {
		db = db.Where("category && ?", info.Category)
	}
	// 根据媒体类型过滤
	if info.DocType != "" {
		if info.DocType == "video" {
			db = db.Where("doc_type = 'video'")
		} else {
			db = db.Where("doc_type != 'video'")
		}

	}
	var tutorialList []model.Tutorial
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&tutorialList).Error
	return tutorialList, total, err
}

func CreateTutorial(tutorial model.Tutorial) (res model.Tutorial, err error) {
	err = global.DB.Create(&tutorial).Error
	// 打包
	go Pack(request.PackRequest{ID: tutorial.ID})
	return tutorial, err
}

func GetTutorial(id uint) (result model.Tutorial, err error) {
	db := global.DB.Model(&model.Tutorial{})
	err = db.Where("id = ?", id).First(&result).Error
	return result, err
}

func DeleteTutorial(req request.DelTutorialRequest) (err error) {
	err = global.DB.Where("id = ?", req.Id).Delete(&model.Tutorial{}).Error
	return err
}

func UpdateTutorial(tutorial model.Tutorial) (err error) {
	raw := global.DB.Where("id = ?", tutorial.ID).Updates(&tutorial)
	if raw.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	go Pack(request.PackRequest{ID: tutorial.ID})
	return raw.Error
}

func UpdateTutorialStatus(id uint, status uint8) (err error) {
	raw := global.DB.Model(&model.Tutorial{}).Where("id = ?", id).Update("status", status)
	if raw.RowsAffected == 0 {
		return errors.New("修改失败")
	}
	return raw.Error
}
