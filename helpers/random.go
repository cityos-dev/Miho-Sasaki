package helpers

import (
	"math/rand"
	"time"
)

func MakeRandomStr(digits int) (string, error) {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, digits)
	for i := 0; i < digits; i++ {
		n := i % 3
		if n < 2 {
			// 65 => 'A', 90 => 'Z'
			bytes[i] = byte(randInt(65, 90))
		} else {
			// 65 => 'a', 90 => 'z'
			bytes[i] = byte(randInt(97, 122))
		}

	}
	return string(bytes), nil
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
