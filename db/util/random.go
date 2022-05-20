package util

import (
	"strings"
	"time"

	"math/rand"
)

const alphabet = "abcdefghiklmnopkrstuvwxyz"

//nit would be called automatically when the package is first used
func init() {

	//the seed  func would make sure everytime the code is run, the generated value would be different..
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generate a random integer Between min and  max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0 -> max-min
}

func Randomstring(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//RandomOwner generate a random owner name
func RandomOwner() string {
	return Randomstring(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generate a  random currency code
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}

	n := len(currencies) // comput the list of the currencies and assign it to n
	return currencies[rand.Intn(n)]
}
