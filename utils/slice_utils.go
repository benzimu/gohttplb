package utils

import (
	"strings"
)

// AddSchemeSlice add http scheme for string slice
func AddSchemeSlice(sli []string) []string {
	result := make([]string, len(sli))
	for i, s := range sli {
		if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
			continue
		}
		s = "http://" + s
		result[i] = s
	}
	return result
}

// TrimStringSlice trim space for string slice
func TrimStringSlice(sli []string) []string {
	result := make([]string, 0, len(sli))
	for _, s := range sli {
		sTmp := strings.TrimSpace(s)
		if sTmp == "" {
			continue
		}
		result = append(result, sTmp)
	}
	return result
}

// ExistStringSlice check exist for string slice
func ExistStringSlice(s string, sli []string) bool {
	for _, tmp := range sli {
		if tmp == s {
			return true
		}
	}
	return false
}

// ExistIntSlice check exist for int slice
func ExistIntSlice(s int, sli []int) bool {
	for _, tmp := range sli {
		if tmp == s {
			return true
		}
	}
	return false
}

// RemoveDuplicateElement remove duplicate element for string slice
func RemoveDuplicateElement(elements []string) []string {
	result := make([]string, 0, len(elements))
	temp := map[string]struct{}{}
	for _, item := range elements {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
