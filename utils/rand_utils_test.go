package utils

import "testing"

func TestGenRandIntn(t *testing.T) {
	for index := 0; index < 5; index++ {
		n := GenRandIntn()
		t.Log("no param: ", n)
	}

	for index := 0; index < 5; index++ {
		n := GenRandIntn(100)
		t.Log("one param: ", n)
	}

	for index := 0; index < 5; index++ {
		n := GenRandIntn(100, 500)
		t.Log("two param: ", n)
	}

	t.Logf("return 0 == %d", GenRandIntn(500, 100))

	t.Error("panic", GenRandIntn(0))
}
