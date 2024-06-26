package model

type CollectionRelate struct {
	ID           uint   `gorm:"primarykey"`
	CollectionID uint   `gorm:"column:collection_id;not null;comment:合辑ID;index:collection_quest_id,UNIQUE" json:"collection_id"`
	QuestID      uint   `gorm:"column:quest_id;not null;comment:题目ID;index:collection_quest_id,UNIQUE" json:"quest_id"`
	TokenID      string `gorm:"column:token_id;not null;comment:token_id;type:varchar(100)" json:"token_id"`
	Sort         int    `gorm:"column:sort;type:int;default:0;comment:排序" json:"sort"`
	Status       uint8  `gorm:"column:status;default:1" json:"status"` // 状态 1 上架 2 未上架
}

func (CollectionRelate) TableName() string {
	return "collection_relate"
}
