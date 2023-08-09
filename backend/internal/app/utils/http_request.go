package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// GetRequest 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func GetRequest(url string) (res string, err error) {
	tryTimes := 3 // 重试次数
	for i := 0; i < tryTimes; i++ {
		// 超时时间：60秒
		client := &http.Client{Timeout: time.Second * 60}
		req, errReq := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", "application/json")
		if errReq != nil {
			err = errReq
			continue
		}
		resp, errReq := client.Do(req)
		if errReq != nil {
			err = errReq
			resp.Body.Close()
			continue
		}
		result, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return string(result), nil
	}
	return "", err
}

// PostRequest 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func PostRequest(url string, data interface{}, contentType string) (res []byte, err error) {
	tryTimes := 3 // 重试次数
	for i := 0; i < tryTimes; i++ {
		// 超时时间：60秒
		client := &http.Client{Timeout: 60 * time.Second}
		jsonStr, _ := json.Marshal(data)
		resp, errReq := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
		if errReq != nil {
			err = errReq
			resp.Body.Close()
			continue
		}
		result, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return result, nil
	}
	return res, err
}

// PostFileRequest 发送POST请求 上传文件
// url：         请求地址
// data：        POST请求提交的数据
func PostFileRequest(url string, header *multipart.FileHeader) (res []byte, err error) {
	file, err := header.Open()
	if err != nil {
		return res, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	//创建一个multipart类型的写文件
	writer := multipart.NewWriter(body)
	//使用给出的属性名paramName和文件名filePath创建一个新的form-data头
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		fmt.Println(" post err=", err)
	}
	//将源复制到目标，将file写入到part   是按默认的缓冲区32k循环操作的，不会将内容一次性全写入内存中,这样就能解决大文件的问题
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		fmt.Println(" post err=", err)
	}
	tryTimes := 3 // 重试次数
	for i := 0; i < tryTimes; i++ {
		request, err := http.NewRequest("POST", url, body)
		if err != nil {
			return nil, err
		}
		//writer.FormDataContentType() ： 返回w对应的HTTP multipart请求的Content-Type的值，多以multipart/form-data起始
		request.Header.Set("Content-Type", writer.FormDataContentType())
		// 超时时间：60秒
		client := &http.Client{Timeout: 60 * time.Second}
		result, errReq := client.Do(request)
		if errReq != nil {
			err = errReq
			//http返回的response的body必须close,否则就会有内存泄露
			result.Body.Close()
			continue
		}
		res, _ = io.ReadAll(result.Body)
		result.Body.Close()
		return res, nil
	}
	return res, err
}

func GraphQlRequest(url string, payload *strings.Reader) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
