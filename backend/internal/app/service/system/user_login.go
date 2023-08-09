package system

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
)

func CreateUserLogin(UserLogin model.UserLogin) (err error) {
	err = global.DB.Create(&UserLogin).Error
	return err
}
