package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strconv"
)

// GetQuestList 获取挑战列表
func GetQuestList(req request.GetQuestListRequest) (res []response.GetQuestListRes, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.DB.Model(&model.Quest{}).Where("style = 1")
	db.Where("disabled = false")
	db.Where(&req.Quest)
	db.Order("sort desc,token_id desc")
	if req.OrderKey == "token_id" {
		fmt.Println(req.OrderKey)
		fmt.Println(req.Desc)
		if req.Desc {
			db.Order("token_id desc")
		} else {
			db.Order("token_id asc")
		}
	} else {
		db.Order("token_id desc")
	}
	if req.SearchKey != "" {
		db.Where("quest.title ILIKE ? OR quest.description ILIKE ?", "%"+req.SearchKey+"%", "%"+req.SearchKey+"%")
		tokenID, err := strconv.Atoi(req.SearchKey)
		if err == nil {
			db.Or("quest.token_id = ?", tokenID)
		}
	}
	if req.Address != "" {
		db.Select("quest.*,c.claimed")
		db.Joins("LEFT JOIN user_challenges c ON quest.token_id = c.token_id AND c.address = ?", req.Address)
	} else {
		db.Select("*")
	}
	err = db.Count(&total).Error
	if err != nil {
		return res, total, err
	}
	err = db.Limit(limit).Offset(offset).Find(&res).Error
	for i := 0; i < len(res); i++ {
		// 统计铸造数量
		global.DB.Model(&model.UserChallenges{}).Where("token_id = ?", res[i].TokenId).Count(&res[i].ClaimNum)
		// 统计挑战人次
		global.DB.Model(&model.UserChallengeLog{}).Where("token_id = ?", res[i].TokenId).Count(&res[i].ChallengeNum)
		// 获取挑战合辑
		global.DB.Model(&model.CollectionRelate{}).Select("collection_id").Where("token_id = ?", res[i].TokenId).Find(&res[i].CollectionID)
	}

	return res, total, err
}

// GetQuest 获取挑战详情
func GetQuest(id string) (quest response.GetQuestRes, err error) {
	err = global.DB.Model(&model.Quest{}).Where("token_id", id).First(&quest).Error
	// 获取所属合辑
	global.DB.Model(&model.CollectionRelate{}).Select("collection_id").Where("token_id = ?", id).Find(&quest.CollectionID)
	return
}

// TopQuest 置顶挑战
func TopQuest(req request.TopQuestRequest) error {
	for _, id := range req.ID {
		err := global.DB.Model(&model.Quest{}).Where("id = ?", id).Update("sort", math.MaxInt).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateQuestStatus 更改上架状态
func UpdateQuestStatus(req request.UpdateQuestStatusRequest) error {
	err := global.DB.Model(&model.Quest{}).Where("id = ?", req.ID).Update("status", req.Status).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateQuest 修改挑战
func UpdateQuest(req request.UpdateQuestRequest) error {
	if req.CollectionID == nil {
		return errors.New("参数错误")
	}
	data := map[string]interface{}{}
	if req.EstimateTime != nil {
		data["quest_data"] = gorm.Expr(fmt.Sprintf("jsonb_set(quest_data, '{estimateTime}', '%d')", *req.EstimateTime))
	} else {
		data["quest_data"] = gorm.Expr(fmt.Sprintf("jsonb_set(quest_data, '{estimateTime}', 'null')"))
	}
	if req.Difficulty != nil {
		data["meta_data"] = gorm.Expr(fmt.Sprintf("jsonb_set(meta_data, '{attributes,difficulty}', '%d')", *req.Difficulty))
	} else {
		data["meta_data"] = gorm.Expr(fmt.Sprintf("jsonb_set(meta_data, '{attributes,difficulty}', 'null')"))
	}
	if req.Sort != nil {
		data["sort"] = *req.Sort
	}

	tx := global.DB.Begin()
	// 查询quest
	var quest model.Quest
	err := tx.Model(&model.Quest{}).Where("id = ?", req.ID).First(&quest).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 清除原有关系
	err = tx.Model(&model.CollectionRelate{}).Where("quest_id = ?", req.ID).Delete(&model.CollectionRelate{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 写入collection关系表
	var exist bool
	for _, id := range *req.CollectionID {
		// 判断集合是否存在
		var collection model.Collection
		err = tx.Model(&model.Collection{}).Where("id = ?", id).First(&collection).Error
		if err != nil {
			tx.Rollback()
			return errors.New("集合不存在")
		}
		// 写入关系
		err = tx.Model(&model.CollectionRelate{}).Create(&model.CollectionRelate{
			CollectionID: id,
			QuestID:      req.ID,
			TokenID:      quest.TokenId,
			Status:       collection.Status,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if exist == false && collection.Status == 1 {
			exist = true
		}
	}
	if len(*req.CollectionID) != 0 && exist {
		data["collection_status"] = 2
	} else {
		data["collection_status"] = 1
	}
	// 更新Quest
	raw := tx.Model(&model.Quest{}).Where("id = ?", req.ID).Updates(data)
	if raw.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	if raw.Error != nil {
		tx.Rollback()
		return raw.Error
	}
	UpdateCollectionStatusAuto(tx) // 更新合辑下架状态
	return tx.Commit().Error
}

// DeleteQuest 删除挑战
func DeleteQuest(req request.DeleteQuestRequest) error {
	raw := global.DB.Model(&model.Quest{}).Where("id = ?", req.ID).Update("disabled", true)
	if raw.RowsAffected == 0 {
		return errors.New("删除失败")
	}
	UpdateCollectionStatusAuto(global.DB) // 更新合辑下架状态
	return raw.Error
}

// UpdateCollectionStatusAuto 更新合辑下架状态
func UpdateCollectionStatusAuto(tx *gorm.DB) error {
	var collectionList []model.Collection
	err := tx.Model(&model.Collection{}).Find(&collectionList).Error
	if err != nil {
		return err
	}
	for _, v := range collectionList {
		var count int64
		err = tx.Model(&model.CollectionRelate{}).Where("collection_id = ?", v.ID).Where("status = 1").Count(&count).Error
		if err != nil {
			return err
		}
		// 更改下架状态
		if count == 0 {
			err = tx.Model(&model.Collection{}).Where("id = ?", v.ID).Update("status", 2).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
