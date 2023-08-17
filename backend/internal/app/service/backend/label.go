package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"errors"
)

// LabelLangList 获取语言列表
func LabelLangList() (language []model.Language, err error) {
	db := global.DB.Model(&model.Language{})
	err = db.Find(&language).Error
	return
}

// LabelAddLang 添加语言
func LabelAddLang(data model.Language) error {
	db := global.DB.Model(&model.Language{})
	return db.Create(&data).Error
}

// LabelRemoveLang 删除语言
func LabelRemoveLang(data model.Language) error {
	db := global.DB.Model(&model.Language{})
	raw := db.Where("id = ?", data.ID).Delete(&model.Language{})
	if raw.RowsAffected == 0 {
		return errors.New("删除失败")
	}
	return raw.Error
}

// LabelUpdateLang 修改语言
func LabelUpdateLang(data model.Language) error {
	db := global.DB.Model(&model.Language{})
	return db.Where("id = ?", data.ID).Updates(&data).Error
}

// LabelCategoryList 获取分类标签列表
func LabelCategoryList() (category []model.Category, err error) {
	db := global.DB.Model(&model.Category{})
	err = db.Find(&category).Error
	return
}

// LabelAddCategory 添加分类标签
func LabelAddCategory(data model.Category) error {
	db := global.DB.Model(&model.Category{})
	return db.Create(&data).Error
}

// LabelRemoveCategory 删除分类标签
func LabelRemoveCategory(data model.Category) error {
	db := global.DB.Model(&model.Category{})
	raw := db.Where("id = ?", data.ID).Delete(&model.Category{})
	if raw.RowsAffected == 0 {
		return errors.New("删除失败")
	}
	return raw.Error
}

// LabelUpdateCategory 修改分类标签
func LabelUpdateCategory(data model.Category) error {
	db := global.DB.Model(&model.Category{})
	return db.Where("id = ?", data.ID).Updates(&data).Error
}

// LabelThemeList 获取分类标签列表
func LabelThemeList() (theme []model.Theme, err error) {
	db := global.DB.Model(&model.Theme{})
	err = db.Find(&theme).Error
	return
}

// LabelAddTheme 添加主题标签
func LabelAddTheme(data model.Theme) error {
	db := global.DB.Model(&model.Theme{})
	return db.Create(&data).Error
}

// LabelRemoveTheme 删除主题标签
func LabelRemoveTheme(data model.Theme) error {
	db := global.DB.Model(&model.Theme{})
	raw := db.Where("id = ?", data.ID).Delete(&model.Theme{})
	if raw.RowsAffected == 0 {
		return errors.New("删除失败")
	}
	return raw.Error
}

// LabelUpdateTheme 修改主题标签
func LabelUpdateTheme(data model.Theme) error {
	db := global.DB.Model(&model.Theme{})
	return db.Where("id = ?", data.ID).Updates(&data).Error
}
