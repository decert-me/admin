package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/response"
	"backend/internal/app/utils/upload"
	"mime/multipart"
	"strings"
)

// UploadAvatar 上传头像
func UploadAvatar(userID uint, header *multipart.FileHeader) (err error, file response.FileUploadResponse) {
	return UploadFile(userID, header, "avatar")
}

func UploadFile(userID uint, header *multipart.FileHeader, _type string) (err error, file response.FileUploadResponse) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		return
	}
	// 保存上传文件记录到数据库
	s := strings.Split(header.Filename, ".")
	fileUpload := model.Upload{
		UserID: userID,
		Url:    filePath,
		Name:   header.Filename,
		Tag:    s[len(s)-1],
		Key:    key,
		Type:   _type,
	}
	if err := global.DB.Model(&model.Upload{}).Create(&fileUpload).Error; err != nil {
		return err, file
	}
	// 返回结果
	f := response.FileUploadResponse{
		Url:  filePath,
		Name: header.Filename,
		Tag:  s[len(s)-1],
		Key:  key,
	}
	return err, f
}
