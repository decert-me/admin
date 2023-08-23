package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// GetQuestList 获取挑战列表
func GetQuestList(req request.GetQuestListRequest) (res []response.GetQuestListRes, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.DB.Model(&model.Quest{})

	db.Where(&req.Quest)
	err = db.Count(&total).Error
	if err != nil {
		return res, total, err
	}
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
	}
	if req.Address != "" {
		db.Select("quest.*,c.claimed")
		db.Joins("LEFT JOIN user_challenges c ON quest.token_id = c.token_id AND c.address = ?", req.Address)
	} else {
		db.Select("*")
	}
	err = db.Limit(limit).Offset(offset).Find(&res).Error
	for i := 0; i < len(res); i++ {
		// 统计铸造数量
		global.DB.Model(&model.UserChallenges{}).Where("token_id = ?", res[i].TokenId).Count(&res[i].ClaimNum)
		// 统计挑战人次
		global.DB.Model(&model.UserChallengeLog{}).Where("token_id = ?", res[i].TokenId).Count(&res[i].ChallengeNum)
	}

	return res, total, err
}

// TopQuest 置顶挑战
func TopQuest(req request.TopQuestRequest) error {
	for i, id := range req.ID {
		err := global.DB.Model(&model.Quest{}).Where("token_id = ?", id).Update("top", req.Top[i]).Error
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
	// 若此挑战已经有人挑战过，不可编辑
	var count int64
	err := global.DB.Model(&model.UserChallenges{}).Where("token_id = ?", req.ID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("此挑战已经有人挑战过，不可编辑")
	}
	data := map[string]interface{}{
		"meta_data":  gorm.Expr(fmt.Sprintf("jsonb_set(data, '{attributes,difficulty}', '%d')", req.Difficulty)),
		"quest_data": gorm.Expr(fmt.Sprintf("jsonb_set(data, '{estimateTime}', '%d')", req.EstimateTime)),
	}
	raw := global.DB.Model(&model.Quest{}).Where("id = ?", req.ID).Updates(data)
	if raw.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	if raw.Error != nil {
		return raw.Error
	}
	return nil
}
