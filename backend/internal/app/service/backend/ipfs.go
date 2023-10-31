package backend

import (
	"backend/internal/app/global"
	"backend/internal/app/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"strconv"
	"sync/atomic"
	"time"
)

type ipfsRPC struct {
	index int
	lock  atomic.Bool
}

var ipfsPoint ipfsRPC

func GetIPFSUploadAPI() string {
	return global.CONFIG.IPFS[ipfsPoint.index].UploadAPI
}

func BalanceIPFS() {
	if ipfsPoint.lock.Load() {
		return
	}
	ipfsPoint.lock.Store(true)
	defer ipfsPoint.lock.Store(false)

	IPFS := global.CONFIG.IPFS
	indexList := make([]int64, len(IPFS))
	for i, v := range IPFS {
		if v.API == "" || v.UploadAPI == "" {
			return
		}
		spent, err := ipfsRequest(v.API, v.UploadAPI)
		if err != nil {
			fmt.Println(err)
		}
		indexList[i] = spent
		time.Sleep(time.Second * 1)
	}
	fmt.Println(indexList)
	ipfsPoint.index, _ = utils.SliceMin[int64](indexList)
	global.LOG.Info("IPFS 切换: " + strconv.Itoa(ipfsPoint.index))
}

func ipfsRequest(api string, uploadAPI string) (spent int64, err error) {
	max := int64(9999999999999)
	defer func() {
		if err := recover(); err != nil {
			spent = max
			return
		}
	}()
	client := req.C().SetTimeout(15 * time.Second)
	startTime := time.Now()
	// 上传JSON
	// 组成请求体
	jsonReq := make(map[string]interface{})
	jsonReq["body"] = "{\"foo\":\"bar\"}"
	// 发送请求
	url := fmt.Sprintf("%s/upload/json", uploadAPI)
	res, err := client.R().SetBody(jsonReq).Post(url)
	if err != nil {
		return max, err
	}
	// 解析返回结果
	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Hash    string `gorm:"column:hash" json:"hash" form:"hash"`
	}
	var resJson Response
	err = json.Unmarshal(res.Bytes(), &resJson)
	if err != nil {
		return max, err
	}
	if resJson.Status != "1" {
		return max, err
	}
	// 请求JSON
	urlReq := fmt.Sprintf("%s/%s", api, resJson.Hash)
	content, err := client.R().Get(urlReq)
	if err != nil || !gjson.Valid(content.String()) {
		return max, err
	}
	return time.Since(startTime).Milliseconds(), nil
}

// IPFSUploadFile
// @description: 上传文件
// @param: header *multipart.FileHeader
// @return: err error, list interface{}, total int64
func IPFSUploadFile(header *multipart.FileHeader) (err error, hash string) {
	file, err := header.Open()
	if err != nil {
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return
	}
	// Convert the byte slice to an io.Reader
	reader := bytes.NewReader(content)

	// 发送请求
	url := fmt.Sprintf("%s/upload/image", GetIPFSUploadAPI())
	client := req.C().SetTimeout(120 * time.Second)
	res, err := client.R().SetFileUpload(req.FileUpload{
		ParamName: "file",
		FileName:  header.Filename,
		GetFileContent: func() (io.ReadCloser, error) {
			return io.NopCloser(reader), nil
		},
		ContentType: header.Header.Get("Content-Type"),
	}).Post(url)
	if err != nil {
		go BalanceIPFS()
		return err, hash
	}
	// 解析返回结果
	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Hash    string `gorm:"column:hash" json:"hash" form:"hash"`
	}
	var resJson Response
	err = json.Unmarshal(res.Bytes(), &resJson)
	if err != nil {
		return err, hash
	}
	if resJson.Status != "1" {
		go BalanceIPFS()
		global.LOG.Error("upload file failed", zap.Error(err))
		return err, hash
	}
	return err, resJson.Hash
}
