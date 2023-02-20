package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
