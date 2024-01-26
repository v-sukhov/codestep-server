package services

import (
	"math/rand"
	"strconv"
)

func generatePassword() string {
	N := rand.Intn(900000) + 100000
	return strconv.Itoa(N)
}
