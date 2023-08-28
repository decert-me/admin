package model

import "gorm.io/gorm"

// Collection 合辑
type Collection struct {
	gorm.Model
	Title       string `gorm:"column:title;not null;comment:合辑标题" json:"title"`
	Description string `gorm:"column:description;comment:合辑简介" json:"description"`
	Cover       string `gorm:"column:cover;comment:封面图" json:"cover"`
	Author      string `gorm:"column:author;type:varchar(64);not null;comment:合辑作者" json:"author"`
	Difficulty  uint8  `gorm:"column:difficulty;type:int2;not null;comment:难度" json:"difficulty"` //0:easy;1:moderate;2:difficult
}
