package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuioplkjhgfdsazxcvbnm"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(m int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < m; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomAmount generates a random amount of Amount
func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

// RandomUnit generates a random Unit name
func RandomUnit() string {
	units := []string{KWH}
	n := len(units)
	return units[rand.Intn(n)]
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
