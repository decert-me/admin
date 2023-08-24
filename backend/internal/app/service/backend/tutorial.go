package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"errors"
	"gorm.io/gorm"
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
	err = db.Limit(limit).Offset(offset).Order("top desc,created_at desc").Find(&tutorialList).Error
	return tutorialList, total, err
}

func CreateTutorial(tutorial model.Tutorial) (res model.Tutorial, err error) {
	err = global.DB.Create(&tutorial).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return res, errors.New("目录名重复")
		}
		return res, err
	}
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
	Pack(request.PackRequest{ID: tutorial.ID})
	return raw.Error
}

func UpdateTutorialStatus(id uint, status uint8) (err error) {
	raw := global.DB.Model(&model.Tutorial{}).Where("id = ? AND pack_status=2", id).Update("status", status)
	if raw.RowsAffected == 0 {
		return errors.New("上架失败，请查看打包状态")
	}
	if raw.Error != nil {
		if raw.Error == gorm.ErrDuplicatedKey {
			return errors.New("目录名重复")
		}
		return raw.Error
	}
	return raw.Error
}

func TopTutorial(req request.TopTutorialRequest) (err error) {
	for i := 0; i < len(req.ID); i++ {
		err = global.DB.Model(&model.Tutorial{}).Where("id = ?", req.ID[i]).Update("top", req.Top[i]).Error
		if err != nil {
			return
		}
	}

	return nil
}
