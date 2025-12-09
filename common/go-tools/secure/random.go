package secure

import (
	"math/rand"
	"time"
)

const (
	// RuneAlNumCS enumerates Alphanumeric case sensitive runes
	RuneAlNumCS = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`

	// RuneAlNumCI enumerates Alphanumeric case insensitive runes
	RuneAlNumCI = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`

	// RuneNumeric numbers only
	RuneNumeric = `0123456789`

	// RuneAlpha enumerates Alpha case sensitive runes
	RuneAlpha = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
)

// RandomString generates a random string based on runes
func RandomString(max int, runes []rune) string {
	rand.Seed(time.Now().UTC().UnixNano())

	b := make([]rune, max)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return string(b)
}
