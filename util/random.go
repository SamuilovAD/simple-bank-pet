package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(length int) string {
	var stringBuilder strings.Builder
	alphabetLength := len(alphabet)
	for i := 0; i < length; i++ {
		symbol := alphabet[rand.Intn(alphabetLength)]
		stringBuilder.WriteByte(symbol)
	}

	return stringBuilder.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(1, 1000)
}

func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}
