package gohttplb

import (
	"math/rand"
	"strings"
	"time"
)

// GenRandIntn return rand int
func GenRandIntn(n ...int) int {
	rand.Seed(time.Now().UnixNano())
	if len(n) == 0 {
		return rand.Int()
	} else if len(n) == 1 {
		return rand.Intn(n[0])
	} else if len(n) == 2 && n[0] < n[1] {
		return rand.Intn(n[1]-n[0]) + n[0]
	}
	return 0
}

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
