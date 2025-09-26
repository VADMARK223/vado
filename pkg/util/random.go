package util

import (
	"math/rand"
	"time"
)

// Один генератор с уникальным сидом
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RndBool() bool {
	return r.Intn(2) == 1
}

func RndIntn(n int) int {
	return r.Intn(n)
}
