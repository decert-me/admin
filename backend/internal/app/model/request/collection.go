package request

import "backend/internal/app/model"

type CreateCollectionRequest struct {
	model.Collection
}

type GetCollectionListRequest struct {
	PageInfo
}

type GetCollectionDetailRequest struct {
	ID uint `json:"id"`
}

type UpdateCollectionRequest struct {
	model.Collection
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
	ID []uint `json:"id"`
	//CollectionSort []int  `json:"collection_sort"`
}
