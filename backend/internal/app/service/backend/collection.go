package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model/request"
)

func CreateCollection(r request.CreateCollectionRequest) error {
	return global.DB.Create(&r.Collection).Error
}
