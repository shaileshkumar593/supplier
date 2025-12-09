package secure

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
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

// HMAC calculates the hash using HMAC (SHA256)
func HMAC(data []byte, key []byte) []byte {
	m := hmac.New(sha256.New, key)
	m.Write(data)
	return m.Sum(nil)
}

// MD5 calculates the MD5 hash of a string
func MD5(raw string) string {
	m := md5.New()
	io.WriteString(m, raw)
	return fmt.Sprintf("%x", m.Sum(nil))
}

// SHA256 calculates the 256-bit (32bytes) secure hash algorithm (SHA) of a string
func SHA256(raw string) string {
	m := sha256.New()
	io.WriteString(m, raw)
	return fmt.Sprintf("%x", m.Sum(nil))
}
