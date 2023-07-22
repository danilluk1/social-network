package utilities

import (
	"math/rand"
)

func GenerateEasyCode() int {
	easyDigits := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	code := 0
	for i := 0; i < 4; i++ {
		digit := easyDigits[rand.Intn(len(easyDigits))]
		code = code*10 + digit
	}

	return code
}
