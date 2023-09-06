package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/receive"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"os"
	"path"
	"sync"
	"time"
)

var l sync.Mutex

// GetPackList 获取打包列表
func GetPackList(info request.PageInfo) (list []response.PackListResponse, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&model.Tutorial{})
	db.Where("pack_status != 1")
	var tutorialList []model.Tutorial
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&tutorialList).Error
	for _, v := range tutorialList {
		list = append(list, response.PackListResponse{
			ID:         v.ID,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
			Img:        v.Img,
			Label:      v.Label,
			PackStatus: v.PackStatus,
			Branch:     v.Branch,
			DocPath:    v.DocPath,
			CommitHash: v.CommitHash,
		})
	}
	return list, total, err
}

// GetPackLog 获取打包日志
func GetPackLog(req request.GetPackLogRequest) (packLog []model.PackLog, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DB.Model(&model.PackLog{})
	db.Where("tutorial_id = ?", req.ID)
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&packLog).Error
	return packLog, total, err
}

// Pack  打包
func Pack(r request.PackRequest) error {
	l.Lock()
	defer l.Unlock()
	// 查询数据库
	var tutorial model.Tutorial
	err := global.DB.Model(&model.Tutorial{}).Where("id = ?", r.ID).First(&tutorial).Error
	if err != nil {
		UpdateTutorialPackStatus(r.ID, 3)
		return err
	}
	client := req.C().SetTimeout(10 * time.Minute)
	if global.CONFIG.Pack.Server == "" {
		UpdateTutorialPackStatus(r.ID, 3)
		return errors.New("打包服务器地址不能为空")
	}
	// 发送请求
	url := fmt.Sprintf("%s/pack/pack", global.CONFIG.Pack.Server)
	res, err := client.R().SetBodyJsonMarshal(tutorial).Post(url)
	if err != nil {
		UpdateTutorialPackStatus(r.ID, 3)
		return err
	}
	// 解析JSON
	var packReceive receive.PackReceive
	err = json.Unmarshal(res.Bytes(), &packReceive)
	if err != nil {
		UpdateTutorialPackStatus(r.ID, 3)
		return err
	}
	if packReceive.Code != 0 {
		UpdateTutorialPackStatus(r.ID, 3)
		return errors.New(gjson.Get(res.String(), "msg").String())
	}

	// 写入打包日志
	err = global.DB.Model(&model.PackLog{}).Create(&model.PackLog{
		TutorialID: r.ID,
		Status:     packReceive.Data.PackLog.Status,
	}).Error
	if err != nil {
		return err
	}
	// 写入日志
	err = global.DB.Model(&model.Tutorial{}).Where("id = ?", r.ID).
		Updates(&packReceive.Data.Tutorial).Error
	if err != nil {
		return err
	}
	// 下载文件
	downUrl := fmt.Sprintf("%s/resource/%s", global.CONFIG.Pack.Server, packReceive.Data.FileName)
	fileRes, err := client.R().Get(downUrl)
	if err != nil {
		UpdateTutorialPackStatus(r.ID, 3)
		return err
	}
	// 写入文件
	file, err := os.Create(fmt.Sprintf("%s/%s.zip", global.CONFIG.Pack.PublishPath, tutorial.CatalogueName))
	if err != nil {
		UpdateTutorialPackStatus(r.ID, 3)
		global.LOG.Error("创建文件失败", zap.Error(err))
		return err
	}
	_, err = file.Write(fileRes.Bytes())
	if err != nil {
		global.LOG.Error("写入文件失败", zap.Error(err))
		UpdateTutorialPackStatus(r.ID, 3)
		return err
	}
	// 解压文件
	zipFilePath := fmt.Sprintf("%s/%s.zip", global.CONFIG.Pack.PublishPath, tutorial.CatalogueName)
	unzipPath := fmt.Sprintf("%s/%s", global.CONFIG.Pack.PublishPath, tutorial.CatalogueName)
	fmt.Println("zipFilePath", zipFilePath)
	fmt.Println("unzipPath", unzipPath)
	err = utils.Unzip(zipFilePath, unzipPath)
	if err != nil {
		global.LOG.Error("解压文件失败", zap.Error(err))
		UpdateTutorialPackStatus(r.ID, 3)
		return err
	}
	// 删除文件
	err = os.Remove(zipFilePath)
	if err != nil {
		global.LOG.Error("删除文件失败", zap.Error(err))
		UpdateTutorialPackStatus(r.ID, 3)
		return err
	}
	return nil
}

func UpdateTutorialPackStatus(id uint, status uint8) (err error) {
	raw := global.DB.Model(&model.Tutorial{}).Where("id = ?", id).Update("pack_status", status)
	if raw.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return raw.Error
}

// PackDelExcessFile 删除多余文件
func PackDelExcessFile(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		// 判断路径是文件还是文件夹
		if !f.IsDir() {
			err = os.Remove(path.Join(dir, f.Name())) // 删除文件
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
