package utils

import (
	"math/rand"
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
