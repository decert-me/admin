package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// CreateCollection 创建合辑
func CreateCollection(r request.CreateCollectionRequest) error {
	collection := model.Collection{
		Title:       r.Title,
		Description: r.Description,
		Cover:       r.Cover,
		Author:      r.Author,
		Style:       2,
		Difficulty:  r.Difficulty,
		Status:      2,
	}
	return global.DB.Model(&model.Collection{}).Create(&collection).Error
}

// GetCollectionList 获取合辑列表
func GetCollectionList(r request.GetCollectionListRequest) (list []response.GetCollectionListRes, total int64, err error) {
	db := global.DB.Model(&model.Collection{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(Paginate(r.Page, r.PageSize)).Order("sort desc,add_ts desc").Find(&list).Error
	for i := 0; i < len(list); i++ {
		// 合辑下的挑战
		var TokenIDList []string
		err := global.DB.Model(&model.CollectionRelate{}).Select("token_id").Where("collection_id = ?", list[i].ID).Find(&TokenIDList).Error
		if err != nil {
			continue
		}
		var claimNumTotal, challengeNumTotal, estimateTimeTotal int64
		for _, tokenId := range TokenIDList {
			var claimNum, challengeNum int64
			// 统计铸造数量
			global.DB.Model(&model.UserChallenges{}).Where("token_id = ?", tokenId).Count(&claimNum)
			// 统计挑战人次
			global.DB.Model(&model.UserChallengeLog{}).Where("token_id = ?", tokenId).Count(&challengeNum)
			// 统计预估时间
			var quest model.Quest
			global.DB.Model(&model.Quest{}).Where("token_id = ?", tokenId).First(&quest)
			estimateTimeTotal += gjson.Get(string(quest.QuestData), "estimateTime").Int()
			claimNumTotal += claimNum
			challengeNumTotal += challengeNum
		}
		list[i].ClaimNum = claimNumTotal
		list[i].ChallengeNum = challengeNumTotal
		list[i].EstimateTime = estimateTimeTotal
		list[i].QuestNum = int64(len(TokenIDList))
	}
	return
}

// GetCollectionDetail 获取合辑详情
func GetCollectionDetail(r request.GetCollectionDetailRequest) (detail model.Collection, err error) {
	err = global.DB.Model(&model.Collection{}).Where("id = ?", r.ID).First(&detail).Error
	return
}

// UpdateCollection 更新合辑
func UpdateCollection(r request.UpdateCollectionRequest) error {
	if r.Sort == nil {
		return errors.New("排序sort不能为空")
	}
	collection := model.Collection{
		Title:       r.Title,
		Description: r.Description,
		Cover:       r.Cover,
		Author:      r.Author,
		Style:       2,
		Sort:        r.Sort,
		Difficulty:  r.Difficulty,
	}
	raw := global.DB.Model(&model.Collection{}).Where("id = ?", r.ID).Updates(&collection)
	if raw.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return raw.Error
}

// DeleteCollection 删除合辑
func DeleteCollection(r request.DeleteCollectionRequest) error {
	tx := global.DB.Begin()
	// 查询合辑状态
	var collection model.Collection
	err := tx.Model(&model.Collection{}).Where("id = ?", r.ID).First(&collection).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 判断合辑是否存在NFT
	if collection.TokenId != "" {
		tx.Rollback()
		return errors.New("合辑已生成NFT，无法删除")
	}
	// 删除合辑
	err = tx.Model(&model.Collection{}).Where("id = ?", r.ID).Delete(&model.Collection{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// CollectionRelate
	var collectionRelateList []model.CollectionRelate
	err = tx.Model(&model.CollectionRelate{}).Where("collection_id = ?", r.ID).Find(&collectionRelateList).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 删除合辑关系表
	err = tx.Model(&model.CollectionRelate{}).Where("collection_id = ?", r.ID).Delete(&model.CollectionRelate{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 判断是否需要更新Quest状态
	for _, v := range collectionRelateList {
		var count int64
		err = tx.Model(&model.CollectionRelate{}).Where("quest_id = ?", v.QuestID).Where("status = 1").Count(&count).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		// 下架
		if count == 0 {
			err = tx.Model(&model.Quest{}).Where("id = ?", v.QuestID).Update("collection_status", 1).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

// UpdateCollectionStatus 更新合辑状态
func UpdateCollectionStatus(r request.UpdateCollectionStatusRequest) error {
	tx := global.DB.Begin()
	err := tx.Model(&model.Collection{}).Where("id = ?", r.ID).Update("status", r.Status).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// CollectionRelate
	var collectionRelateList []model.CollectionRelate
	err = tx.Model(&model.CollectionRelate{}).Where("collection_id = ?", r.ID).Find(&collectionRelateList).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 更新CollectionRelate状态
	err = tx.Model(&model.CollectionRelate{}).Where("collection_id = ?", r.ID).Update("status", r.Status).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 判断是否需要更新Quest状态
	for _, v := range collectionRelateList {
		var count int64
		err = tx.Model(&model.CollectionRelate{}).Where("quest_id = ?", v.QuestID).Where("status = 1").Count(&count).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		fmt.Println("上架状态", r.Status, count)
		// 下架
		if r.Status == 2 && count == 0 {
			err = tx.Model(&model.Quest{}).Where("id = ?", v.QuestID).Update("collection_status", 1).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		// 上架
		if r.Status == 1 && count > 0 {
			err = tx.Model(&model.Quest{}).Where("id = ?", v.QuestID).Update("collection_status", 2).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// GetCollectionQuest 获取合辑下的挑战
func GetCollectionQuest(r request.GetCollectionQuestRequest) (questList []response.GetQuestListRes, err error) {
	err = global.DB.Model(&model.CollectionRelate{}).Select("quest.*").
		Joins("left join quest ON collection_relate.quest_id=quest.id").
		Where("collection_relate.collection_id = ?", r.ID).
		Order("collection_relate.sort desc").Find(&questList).Error
	for i := 0; i < len(questList); i++ {
		// 统计铸造数量
		global.DB.Model(&model.UserChallenges{}).Where("token_id = ?", questList[i].TokenId).Count(&questList[i].ClaimNum)
		// 统计挑战人次
		global.DB.Model(&model.UserChallengeLog{}).Where("token_id = ?", questList[i].TokenId).Count(&questList[i].ChallengeNum)
		// 获取挑战合辑
		global.DB.Model(&model.CollectionRelate{}).Select("collection_id").Where("token_id = ?", questList[i].TokenId).Find(&questList[i].CollectionID)
	}
	return
}

// UpdateCollectionQuestSort 编辑合辑下的挑战排序
func UpdateCollectionQuestSort(r request.UpdateCollectionQuestSortRequest) error {
	for i := 0; i < len(r.ID); i++ {
		err := global.DB.Model(&model.CollectionRelate{}).Where("collection_id = ?", r.CollectionID).Where("quest_id = ?", r.ID[len(r.ID)-i-1]).Update("sort", i).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// AddQuestToCollection 添加挑战到合辑
func AddQuestToCollection(r request.AddQuestToCollectionRequest) error {
	tx := global.DB.Begin()
	// 查询原有Quest ID 列表
	var questIDList []uint
	err := tx.Model(&model.CollectionRelate{}).Select("quest_id").Where("collection_id = ?", r.CollectionID).Find(&questIDList).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 删除原有关系
	//err = tx.Model(&model.CollectionRelate{}).Where("collection_id = ?", r.CollectionID).Delete(&model.CollectionRelate{}).Error
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}
	// 查询合辑状态
	var collection model.Collection
	err = tx.Model(&model.Collection{}).Where("id = ?", r.CollectionID).First(&collection).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 判断合辑是否存在NFT
	if collection.TokenId != "" {
		tx.Rollback()
		return errors.New("合辑已生成NFT，无法修改")
	}
	for _, v := range r.ID {
		// 查询Quest信息
		var quest model.Quest
		err = tx.Model(&model.Quest{}).Where("id = ?", v).First(&quest).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		// 添加到合辑
		collectionRelate := model.CollectionRelate{
			CollectionID: r.CollectionID,
			QuestID:      v,
			TokenID:      quest.TokenId,
			Status:       collection.Status,
		}
		err = tx.Model(&model.CollectionRelate{}).Create(&collectionRelate).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		// 更新Quest状态
		UpdateQuestCollectionStatus(tx, v)
	}
	// 合辑里没有Quest，将Collection下架
	//if len(r.ID) == 0 {
	//	err = tx.Model(&model.Collection{}).Where("id = ?", r.CollectionID).Update("status", 2).Error
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//}
	// 更新Quest状态
	for _, questID := range questIDList {
		UpdateQuestCollectionStatus(tx, questID)
	}
	return tx.Commit().Error
}

// UpdateQuestCollectionStatus 更新Quest是否在合辑状态
func UpdateQuestCollectionStatus(tx *gorm.DB, id uint) {
	var count int64
	err := tx.Model(&model.CollectionRelate{}).Where("quest_id = ?", id).Where("status = 1").Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		// 独立
		tx.Model(&model.Quest{}).Where("id = ?", id).Update("collection_status", 1)
	} else {
		// 合辑
		tx.Model(&model.Quest{}).Where("id = ?", id).Update("collection_status", 2)
	}
}
