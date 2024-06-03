package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/utils"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
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
		// 统计挑战人数
		global.DB.Model(&model.UserChallengeLog{}).Where("token_id = ?", res[i].TokenId).Group("address").Count(&res[i].ChallengeUserNum)
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
	if req.Description != nil {
		// 判断链上是否有数据
		if gjson.Get(string(quest.MetaData), "description").String() != "" {
			tx.Rollback()
			return errors.New("链上已存在数据，无法修改描述")
		}
	}
	if req.Description != nil {
		data["description"] = *req.Description
	}
	// 查询原有关系
	var collectionIDList []uint
	err = tx.Model(&model.CollectionRelate{}).Where("quest_id = ?", req.ID).Pluck("collection_id", &collectionIDList).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 判断 collectionIDList 和 req.CollectionID是否相同
	if utils.CollectionEqual(collectionIDList, *req.CollectionID) {
		// 只更新Quest
		raw := tx.Model(&model.Quest{}).Where("id = ?", req.ID).Updates(data)
		if raw.RowsAffected == 0 {
			tx.Rollback()
			return errors.New("更新失败")
		}
		if raw.Error != nil {
			tx.Rollback()
			return raw.Error
		}
		return tx.Commit().Error
	}

	// 需要移除的关系
	for _, v := range utils.CollectionSubtract(collectionIDList, *req.CollectionID) {
		// 查询合辑状态
		var collection model.Collection
		err := tx.Model(&model.Collection{}).Where("id = ?", v).First(&collection).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		// 判断合辑是否存在NFT
		if collection.TokenId != "" {
			tx.Rollback()
			return errors.New("合辑已生成NFT，无法修改、删除")
		}
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
		// 判断合辑是否存在NFT
		if collection.TokenId != "" {
			tx.Rollback()
			return errors.New("合辑已生成NFT，无法修改、删除")
		}
		var status uint8
		if quest.Status == 2 {
			status = 2
		} else {
			status = collection.Status
		}
		// 写入关系
		err = tx.Model(&model.CollectionRelate{}).Create(&model.CollectionRelate{
			CollectionID: id,
			QuestID:      req.ID,
			TokenID:      quest.TokenId,
			Status:       status,
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
		tx.Rollback()
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

// GetQuestCollectionAddList 获取待添加到合辑挑战列表
func GetQuestCollectionAddList(req request.GetQuestCollectionAddListRequest) (res []response.GetQuestCollectionAddListRes, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DB.Model(&model.Quest{})
	db.Joins("left join collection_relate cr on quest.token_id = cr.token_id")
	db.Where("quest.style = 1")
	db.Where("quest.disabled = false")
	db.Where("cr.id is null")
	if req.SearchKey != "" {
		db.Where("quest.title ILIKE ? OR quest.description ILIKE ?", "%"+req.SearchKey+"%", "%"+req.SearchKey+"%")
		tokenID, err := strconv.Atoi(req.SearchKey)
		if err == nil {
			db.Or("quest.token_id = ?", tokenID)
		}
	}
	db.Order("quest.sort desc,quest.token_id desc")
	err = db.Count(&total).Error
	if err != nil {
		return res, total, err
	}
	err = db.Limit(limit).Offset(offset).Find(&res).Error
	for i := 0; i < len(res); i++ {
		// 获取挑战合辑
		global.DB.Model(&model.CollectionRelate{}).Select("collection_id").Where("token_id = ?", res[i].TokenId).Find(&res[i].CollectionID)
	}
	return
}

// GetQuestStatistics 获取挑战结果详情列表
func GetQuestStatistics(tokenId string) (res []response.GetQuestStatisticsRes, err error) {
	// 查询挑战
	var quest model.Quest
	err = global.DB.Model(&model.Quest{}).Where("token_id", tokenId).First(&quest).Error
	if err != nil {
		return res, err
	}
	// 判断是否是开放题
	if !IsOpenQuest(gjson.Get(string(quest.QuestData), "questions").String()) {
		rankListSQL := `
		WITH ranked AS (
		 SELECT address,ucl.token_id, created_at,ROW_NUMBER() OVER (PARTITION BY address ORDER BY created_at ASC) as rn,
		 (SELECT bool_or(pass) FROM user_challenge_log WHERE address = ucl.address AND token_id = ucl.token_id) as pass,
		 (SELECT max(user_score) FROM user_challenge_log WHERE address = ucl.address AND token_id = ucl.token_id) as highest_score,
		 (SELECT claimed FROM user_challenges WHERE address = ucl.address AND token_id = ucl.token_id) as claimed,
		 EXISTS(SELECT 1 FROM zcloak_card WHERE address = ucl.address AND quest_id = quest.id) as has_did
	     FROM user_challenge_log ucl
			 LEFT JOIN quest ON quest.token_id=ucl.token_id
	     WHERE ucl.token_id = ? AND address !='' AND 			deleted_at IS NULL
		)
		SELECT ROW_NUMBER() OVER (ORDER BY created_at ASC) as rank,ranked.address,pass,ranked.created_at as finish_time,highest_score,	COALESCE(claimed,false) as claimed,has_did
		FROM ranked
		LEFT JOIN users ON ranked.address=users.address
		WHERE rn=1 ORDER BY created_at ASC;
		`
		err = global.DB.Raw(rankListSQL, tokenId).Scan(&res).Error
		if err != nil {
			return res, err
		}
		return res, err
	}
	// 开放题
	rankListSQL := `
		WITH all_open_quest AS(
			SELECT address,token_id,created_at
			FROM user_open_quest
			WHERE token_id = ? AND pass=true
			UNION
			SELECT address,token_id,created_at
			FROM user_challenge_log
			WHERE token_id = ? AND pass=true AND is_open_quest=false
		),ranked_open_quest AS (
		 SELECT address,token_id,created_at,ROW_NUMBER() OVER (PARTITION BY address ORDER BY created_at ASC) as rn 
		 FROM all_open_quest
		 ),
		ranked_with_rank AS (
		 SELECT ROW_NUMBER() OVER (ORDER BY created_at ASC) as rank,address,token_id,created_at as finish_time 
		 FROM ranked_open_quest 
		 WHERE rn=1 
		)
		SELECT ranked_with_rank.*,
		(SELECT bool_or(pass) FROM user_challenge_log WHERE address = ucl.address AND token_id = ucl.token_id) as pass,
		(SELECT MAX(highest_score) 
		FROM (
			(SELECT max(user_score) as highest_score 
			 FROM user_challenge_log 
			 WHERE address = ucl.address 
				 AND token_id = ucl.token_id) 
			UNION ALL 
			(SELECT max(user_score) as highest_score 
			 FROM user_open_quest 
			 WHERE address = ucl.address 
				 AND token_id = ucl.token_id)
		) AS combined_scores) as highest_score,
		 COALESCE((SELECT claimed FROM user_challenges WHERE address = ucl.address AND token_id = ucl.token_id),false) as claimed,
		 EXISTS(SELECT 1 FROM zcloak_card WHERE address = ucl.address AND quest_id = quest.id) as has_did
		FROM ranked_with_rank
		LEFT JOIN user_challenge_log ucl ON ucl.token_id=ranked_with_rank.token_id
		LEFT JOIN quest ON quest.token_id=ucl.token_id
		ORDER BY rank ASC 
		LIMIT 10;
	`
	err = global.DB.Raw(rankListSQL, tokenId, tokenId).Scan(&res).Error
	if err != nil {
		return res, err
	}
	return res, err
}
