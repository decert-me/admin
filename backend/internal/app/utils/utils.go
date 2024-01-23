package utils

import (
	"encoding/json"
	"github.com/google/uuid"
	"reflect"
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

func DiffStructs(s1, s2 interface{}) []string {
	var diff []string

	v1 := reflect.ValueOf(s1)
	v2 := reflect.ValueOf(s2)

	if v1.Kind() != reflect.Struct || v2.Kind() != reflect.Struct {
		return diff
	}

	for i := 0; i < v1.NumField(); i++ {
		field1 := v1.Field(i)
		field2 := v2.Field(i)

		if !reflect.DeepEqual(field1.Interface(), field2.Interface()) {
			fieldName := v1.Type().Field(i).Name
			diff = append(diff, fieldName)
		}
	}

	return diff
}

// IsUUID 判断是否UUID
func IsUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}
