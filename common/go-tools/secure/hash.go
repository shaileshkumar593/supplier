package secure

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/dchest/blake2b"
)

// Blake2B hash/encrypt data using Blake2B algorithm
func Blake2B(data string) string {
	digest, _ := blake2b.New(&blake2b.Config{Size: 8})
	return fmt.Sprintf("%x", digest.Sum([]byte(data)))
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
