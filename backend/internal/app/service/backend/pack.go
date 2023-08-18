package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/model"
	"backend/internal/app/model/request"
	"backend/internal/app/model/response"
	"backend/internal/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"path"
	"strings"
)

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
func Pack(req request.PackRequest) error {
	// 查询数据库
	var tutorial model.Tutorial
	err := global.DB.Model(&model.Tutorial{}).Where("id = ?", req.ID).First(&tutorial).Error
	if err != nil {
		return err
	}
	data := []model.Tutorial{tutorial}
	// 生成JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// 写入JSON
	// 将JSON数据写入文件
	file, err := os.Create(path.Join(global.CONFIG.Pack.Path, "tutorials.json"))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// npm run build -- blockchain-basic
	//args := []string{"run", "build", "--", tutorial.CatalogueName}
	args := []string{"run", "build"}
	dir := global.CONFIG.Pack.Path
	stdoutRes, stdoutErr, err := execCommand(global.CONFIG.Pack.Path, "npm", args...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(stdoutRes)
	fmt.Println(stdoutErr)
	var success bool
	var startPage string
	for _, v := range stdoutRes {
		v = strings.Replace(v, "\n", "", -1)
		// 判断打包是否成功
		if !success && v == "Build completed successfully" {
			success = true
		}
		if gjson.Valid(v) {
			startPage = gjson.Get(v, "startPage").String()
		}
	}
	var status uint8
	if success {
		status = 2
	} else {
		status = 3
	}
	// 写入打包日志
	err = global.DB.Model(&model.PackLog{}).Create(&model.PackLog{
		TutorialID: req.ID,
		Status:     status,
	}).Error
	if err != nil {
		return err
	}
	if status == 2 {
		// 将结果写入数据库
		err = global.DB.Model(&model.Tutorial{}).Where("id = ?", req.ID).
			Updates(&model.Tutorial{StartPage: startPage, PackStatus: status}).Error
		if err != nil {
			return err
		}
	}
	if status == 3 {
		return errors.New("打包失败")
	}
	// 删除多余文件
	//PackDelExcessFile(dir + "/build")
	// 复制文件
	//utils.CopyContents(path.Join(dir, "build", tutorial.CatalogueName), path.Join(dir, "build"))
	// 删除文件
	//os.RemoveAll(path.Join(dir, "build", tutorial.CatalogueName))
	// 复制文件到发布项目路径
	utils.CopyContents(path.Join(dir, "build"), global.CONFIG.Pack.PublishPath)
	// Build completed successfully
	// Error running build command:
	return nil
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

// 打包失败通知
// 打包成功后复制文件

// 解析视频
