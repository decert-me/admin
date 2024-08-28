package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"errors"
	"gorm.io/datatypes"
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

// InitOpenQuestUserScore 初始化用户分数
func InitOpenQuestUserScore() {
	// 初始化用户分数
	var userOpenQuests []model.UserOpenQuest
	err := global.DB.Model(&model.UserOpenQuest{}).Where("user_score IS NULL").Find(&userOpenQuests).Error
	if err != nil {
		return
	}
	for _, userOpenQuest := range userOpenQuests {
		// 获取题目
		var quest model.Quest
		if err = global.DB.Model(&model.Quest{}).Where("token_id = ?", userOpenQuest.TokenId).First(&quest).Error; err != nil {
			continue
		}
		result, err := AnswerCheck(global.CONFIG.Quest.EncryptKey, datatypes.JSON(userOpenQuest.Answer), quest)
		if err != nil {
			continue
		}
		userScore := result.UserScore
		// 更新用户分数
		err = global.DB.Model(&model.UserOpenQuest{}).Where("id = ?", userOpenQuest.ID).Update("user_score", userScore).Error
		if err != nil {
			continue
		}
	}
}
