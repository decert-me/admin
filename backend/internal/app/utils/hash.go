package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashData 生成校验Hash
func HashData(data interface{}, key string) (timestamp int64, hashValue string) {
	// 将结构体数据转换为 JSON 格式的字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	// 获取当前时间戳
	timestamp = time.Now().Unix()

	// 将时间戳与字符串拼接
	hashData := string(jsonData) + key + fmt.Sprintf("%d", timestamp)
	// 对字符串进行哈希计算
	hasher := sha256.New()
	hasher.Write([]byte(hashData))
	hashValue = hex.EncodeToString(hasher.Sum(nil))
	return timestamp, hashValue
}
