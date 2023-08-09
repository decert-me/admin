package utils

import (
	"encoding/json"
)

// SliceIsExist 判断元素是否在slice
func SliceIsExist[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func StructToMap(inter []any) map[string]interface{} {
	var m map[string]interface{}
	for _, v := range inter {
		ja, _ := json.Marshal(v)
		json.Unmarshal(ja, &m)
	}
	return m
}

func MapPushStruct(m map[string]interface{}, inter []any) map[string]interface{} {
	for _, v := range inter {
		ja, _ := json.Marshal(v)
		json.Unmarshal(ja, &m)
	}
	return m
}
