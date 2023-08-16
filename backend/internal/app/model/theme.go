package model

import "gorm.io/gorm"

type Theme struct {
	gorm.Model
	Chinese string
	English string
	Weight  int `gorm:"column:weight;default:0"`
}
