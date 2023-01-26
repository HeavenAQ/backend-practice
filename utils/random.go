package utils

import (
	"math/rand"
	"time"
)

func init() {
	// like srand(time(NULL))
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	str := ""
	for i := 0; i < n; i++ {
		str += string(rune('a' + randomInt(0, 25)))
	}
	return str
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	length := len(currencies)
	return currencies[rand.Intn(length)]
}
