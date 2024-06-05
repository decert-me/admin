package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = -1
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func ResultWithRaw(data []byte, c *gin.Context) {
	c.Data(http.StatusOK, "application/json", data)
	//c.JSON(http.StatusOK, data)
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func FailWithErrorMessage(message string, err error, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, fmt.Sprintf(message+"：%s", err.Error()), c)
}

// TranslateValidationErrors 转换验证错误为中文
func TranslateValidationErrors(err error) string {
	var errMessages []string
	// 断言 error 为 validator.ValidationErrors 类型
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// 如果不是预期的类型，直接返回原始错误信息
		return err.Error()
	}
	// 遍历所有的字段验证错误
	for _, e := range validationErrors {
		switch e.Tag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("%s 不能为空", e.Field()))
		// 可以添加更多的case来处理不同的验证标签
		default:
			errMessages = append(errMessages, fmt.Sprintf("%s 验证失败", e.Field()))
		}
	}

	// 将所有错误信息合并为一个字符串
	return strings.Join(errMessages, "; ")
}
