package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"github.com/tidwall/gjson"
)

// CreateCollection 创建合辑
func CreateCollection(r request.CreateCollectionRequest) error {
	var tokenID int64
	err := global.DB.Model(&model.Quest{}).Select("MIN(token_id)").Find(&tokenID).Error
	if err != nil {
		return err
	}
	collection := model.Quest{
		TokenId:     tokenID - 1,
		Title:       r.Title,
		Difficulty:  r.Difficulty,
		Description: r.Description,
		Cover:       r.Cover,
		Author:      r.Author,
		Style:       2,
	}
	return global.DB.Model(&model.Quest{}).Create(&collection).Error
}

// GetCollectionList 获取合辑列表
func GetCollectionList(r request.GetCollectionListRequest) (list []response.GetCollectionListRes, total int64, err error) {
	db := global.DB.Model(&model.Quest{}).Where("style = 2")
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(Paginate(r.Page, r.PageSize)).Find(&list).Error
	for i := 0; i < len(list); i++ {
		// 合辑下的挑战
		var TokenIDList []int64
		err := global.DB.Model(&model.Quest{}).Select("token_id").Where("collection_id = ?", list[i].ID).Find(&TokenIDList).Error
		if err != nil {
			continue
		}
		var claimNumTotal, challengeNumTotal, estimateTimeTotal int64
		for tokenId := range TokenIDList {
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
	}
	return
}

// GetCollectionDetail 获取合辑详情
func GetCollectionDetail(r request.GetCollectionDetailRequest) (detail model.Quest, err error) {
	err = global.DB.Model(&model.Quest{}).Where("id = ?", r.ID).First(&detail).Error
	return
}

// UpdateCollection 更新合辑
func UpdateCollection(r request.UpdateCollectionRequest) error {
	collection := model.Quest{
		Title:       r.Title,
		Difficulty:  r.Difficulty,
		Description: r.Description,
		Cover:       r.Cover,
		Author:      r.Author,
		Style:       2,
	}
	raw := global.DB.Model(&model.Quest{}).Where("id = ?", r.ID).Updates(&collection)
	if raw.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return raw.Error
}

// DeleteCollection 删除合辑
func DeleteCollection(r request.DeleteCollectionRequest) error {
	return global.DB.Where("id = ?", r.ID).Delete(&model.Quest{}).Error
}

// UpdateCollectionStatus 更新合辑状态
func UpdateCollectionStatus(r request.UpdateCollectionStatusRequest) error {
	return global.DB.Model(&model.Quest{}).Where("id = ?", r.ID).Update("status", r.Status).Error
}

// GetCollectionQuest 获取合辑下的挑战
func GetCollectionQuest(r request.GetCollectionQuestRequest) (list []model.Quest, err error) {
	err = global.DB.Model(&model.Quest{}).Where("collection_id = ?", r.ID).Order("collection_sort desc").Find(&list).Error
	return
}

// UpdateCollectionQuestSort 编辑合辑下的挑战排序
func UpdateCollectionQuestSort(r request.UpdateCollectionQuestSortRequest) error {
	return global.DB.Model(&model.Quest{}).Where("id = ?", r.ID).Update("collection_sort", r.CollectionSort).Error
}