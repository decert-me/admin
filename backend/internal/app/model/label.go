package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	Language pq.StringArray `gorm:"column:language;type:text[];comment:语言" json:"language,omitempty"`
	Category pq.StringArray `gorm:"column:category;type:text[];comment:分类标签" json:"category,omitempty"`
	Theme    pq.StringArray `gorm:"column:theme;type:text[];comment:主题标签" json:"theme,omitempty"`
}
