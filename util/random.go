package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

//ooh, init is called automatically when the function is first
//used
func init() {
	rand.Seed(time.Now().UnixNano())
}

///RandomInt generates a random value between min and max
func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Random Owner generates random name for owner
func RandomOwner() string {
	return RandomString(6)
}

//Generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//Generates a random currency
func RandomCurrency() string {
	currency := []string{EUR, NGN, USD, GBP}
	return currency[rand.Intn(len(currency))]
}
