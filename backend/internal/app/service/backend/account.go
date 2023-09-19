package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
)

// GetAddressInfo 获取地址信息
func GetAddressInfo(req request.GetAddressInfoRequest) (user model.Users, err error) {
	db := global.DB.Model(&model.Users{})
	err = db.Where("address = ?", req.Address).First(&user).Error
	return
}
