package request

type CreateCollectionRequest struct {
	ID          uint   `gorm:"primarykey"`
	AddTs       int64  `gorm:"column:add_ts;autoCreateTime" json:"addTs"`
	Title       string `gorm:"column:title;not null;comment:合辑标题" json:"title"`
	Description string `gorm:"column:description;comment:合辑简介" json:"description"`
	Cover       string `gorm:"column:cover;comment:封面图" json:"cover"`
	Author      string `gorm:"column:author;type:varchar(64);not null;comment:合辑作者" json:"author"`
	Difficulty  *uint8 `gorm:"column:difficulty;type:int2;not null;comment:难度" json:"difficulty"` //0:easy;1:moderate;2:difficult
	Status      *uint8 `gorm:"column:status;type:int2;default:1;comment:上架状态" json:"status"`      // 1:下架;2:上架
	Sort        int    `gorm:"column:sort;type:int;default:0;comment:排序" json:"sort"`
}

type GetCollectionListRequest struct {
	PageInfo
}

type GetCollectionDetailRequest struct {
	ID uint `json:"id"`
}

type UpdateCollectionRequest struct {
	ID          uint   `gorm:"primarykey"`
	AddTs       int64  `gorm:"column:add_ts;autoCreateTime" json:"addTs"`
	Title       string `gorm:"column:title;not null;comment:合辑标题" json:"title"`
	Description string `gorm:"column:description;comment:合辑简介" json:"description"`
	Cover       string `gorm:"column:cover;comment:封面图" json:"cover"`
	Author      string `gorm:"column:author;type:varchar(64);not null;comment:合辑作者" json:"author"`
	Difficulty  *uint8 `gorm:"column:difficulty;type:int2;not null;comment:难度" json:"difficulty"` //0:easy;1:moderate;2:difficult
	Status      *uint8 `gorm:"column:status;type:int2;default:1;comment:上架状态" json:"status"`      // 1:下架;2:上架
	Sort        int    `gorm:"column:sort;type:int;default:0;comment:排序" json:"sort"`
}

type DeleteCollectionRequest struct {
	ID uint `json:"id"`
}

type UpdateCollectionStatusRequest struct {
	ID     uint  `json:"id"`
	Status uint8 `json:"status" binding:"required"`
}

type GetCollectionQuestRequest struct {
	ID uint `json:"id"`
}

type UpdateCollectionQuestSortRequest struct {
	ID           []uint `json:"id"`
	CollectionID uint   `json:"collection_id"`
}

type AddQuestToCollectionRequest struct {
	ID           []uint `json:"id"`
	CollectionID uint   `json:"collection_id"`
}
