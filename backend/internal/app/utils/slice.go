package utils

import "golang.org/x/exp/constraints"

func SliceMin[T constraints.Ordered](slice []T) (index int, m T) {
	for i, e := range slice {
		if i == 0 || e < m {
			m = e
			index = i
		}
	}
	return
}

// CollectionSubtract 合集相减
func CollectionSubtract[T comparable](a, b []T) []T {
	result := make([]T, 0)
	for _, v := range a {
		if !SliceIsExist(b, v) {
			result = append(result, v)
		}
	}
	return result
}

// CollectionEqual 合集相等
func CollectionEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !SliceIsExist(b, v) {
			return false
		}
	}
	return true
}
