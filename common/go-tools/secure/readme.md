# secure
--
    import "bitbucket.org/matchmove/go-tools/secure"


## Usage

```go
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
```

#### func  Blake2B

```go
func Blake2B(data string) string
```
Blake2B hash/encrypt data using Blake2B algorithm

#### func  HMAC

```go
func HMAC(data []byte, key []byte) []byte
```
HMAC calculates the hash using HMAC (SHA256)

#### func  MD5

```go
func MD5(raw string) string
```
MD5 calculates the MD5 hash of a string

#### func  RandomString

```go
func RandomString(max int, runes []rune) string
```
RandomString generates a random string based on runes

#### func  SHA256

```go
func SHA256(raw string) string
```
SHA256 calculates the 256-bit (32bytes) secure hash algorithm (SHA) of a string

#### func  AESEncrypt

```go
func AESEncrypt(text, secretKey string) string
```
AESEncrypt encrypts string using AES with block size of 16, secretKey should be a 32-byte vaue

#### func  AESDecrypt

```go
func AESDecrypt(text, secretKey string) string
```
AESDecrypt decrypts string using AES with block size of 16, secretKey should be a 32-byte vaue

```go
func AESDecryptWithIV(text, secretKey string) (string, error)
```
AESDecrypt decrypts string using AES 256 CBC and make use of initialization vector when decrypting value, useful for encypts via PHP, then decrypts via Go