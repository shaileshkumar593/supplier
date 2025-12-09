package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // JPEG image support
	_ "image/png"  // PNG image support
	"io"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"

	"swallow-supplier/config"
	"swallow-supplier/utils/array"
	"swallow-supplier/utils/constant"
	"swallow-supplier/utils/secure"
)

const (
	// DefaultDateTimeLayout default ISO8601/RFC3339 datetime
	DefaultDateTimeLayout = "2006-01-02T15:04:05-0700"

	// DefaultTimezone default timezone
	DefaultTimezone = "Asia/Singapore"

	typeAlpha   = "alpha"
	typeNumeric = "numeric"

	alphaBytes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numericBytes = "1234567890"
	mixedBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	charset      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// TypeProvider represents code for Provider
	TypeProvider = "P"
	// TypeVerification represents code for Verification
	TypeVerification = "V"
	// TypeReferenceSdk represents code for Sdk
	TypeReferenceSdk = "RS"

	// ConsumerIdentifier resource identifier for creating the uuid for Consumers, e.g csm_<uuid>
	ConsumerIdentifier = "csm"
)

// SliceStringToSliceInt converts a slice of string to slice of ints
func SliceStringToSliceInt(uS []string) (sI []int) {
	for _, item := range uS {
		aI, _ := strconv.Atoi(item)
		sI = append(sI, aI)
	}

	return
}

// UnixTimestamp function - returns utc int timestamp ...
func UnixTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// GenerateUUID generates a version 4 of UUID with option to remove hyphen or not
func GenerateUUID(identifier string, withHyphen bool) string {
	u := uuid.New()

	var uuidValue = u.String()
	if identifier != "" {
		uuidValue = identifier + "_" + uuidValue
	}

	if withHyphen {
		return uuidValue
	}

	return strings.Replace(uuidValue, "-", "", -1)
}

// GenerateRandomHash generates random hashed string
func GenerateRandomHash(hashLength int) string {
	unixTime := UnixTimestamp()
	return fmt.Sprintf("%s%s", unixTime, secure.RandomString(hashLength-len(unixTime), []rune(secure.RuneAlNumCS)))
}

// GenerateRandomMD5 generates random md5 string
func GenerateRandomMD5() string {
	unixTime := UnixTimestamp()
	randString := secure.RandomString(32-len(unixTime), []rune(secure.RuneAlNumCS))
	data := []byte(randString)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// GenerateRandomSHA generates random sha string
func GenerateRandomSHA() string {
	unixTime := UnixTimestamp()
	randString := secure.RandomString(32-len(unixTime), []rune(secure.RuneAlNumCS))
	data := []byte(randString)
	return fmt.Sprintf("%x", sha1.Sum(data))
}

// GenerateRandString will generate a random string
// param n is the number of characters to generate
// param t is the characters to map from ("alpha", "numeric")
func GenerateRandString(n int, t string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		if t == typeAlpha {
			b[i] = alphaBytes[rand.Intn(len(alphaBytes))]
		} else if t == typeNumeric {
			b[i] = numericBytes[rand.Intn(len(numericBytes))]
		} else {
			b[i] = mixedBytes[rand.Intn(len(mixedBytes))]
		}
	}
	return string(b)
}

// GenerateRandAlpha will generate a random string with alphabet characters only
// param n is number of characters to generate
func GenerateRandAlpha(n int) string {
	return GenerateRandString(n, typeAlpha)
}

// GenerateRandNumeric will generate random string with numeric characters only
// param n is the number of characters to generate
func GenerateRandNumeric(n int) string {
	return GenerateRandString(n, typeNumeric)
}

// GenerateRandAlphaNumeric will generate random string with mixed alphabet and numeric characters only
// param n is for the number of characters to generate
func GenerateRandAlphaNumeric(n int) string {
	len := n / 2
	alphaLen := len
	numericLen := n - len

	// string will always be composed on half alphabets and half numbers
	s := GenerateRandAlpha(alphaLen) + GenerateRandNumeric(numericLen)

	return s
}

// GenerateUniqueNumeric will generate unique numeric string based from unix nano time
// param n is number of characters to generate
// Example: 1546510989503519774
func GenerateUniqueNumeric(n int) string {
	t := int64(time.Now().UnixNano())
	s := strconv.FormatInt(t, 10)

	if n > len(s) {
		return s
	}

	// get last digits from timestamp as unique numeric string
	s = s[len(s)-n:]

	return s
}

// GenerateCode will generate unique code
// group :
//
//		P = provider
//		V = verification
//	 RS = reference sdk
func GenerateCode(group string, customRef string, customJoin string) string {
	var (
		s         []string
		microTime = time.Now().UnixNano() / 1e6
	)

	if group == TypeProvider {
		s = []string{
			"M",                             // M for Matchmove
			(config.Instance().AppEnv)[0:2], // first 2 characters of current app env
			"BO",                            // BO for boilerplate
			TypeProvider,                    // P
			customRef,                       // custom reference
			GenerateRandAlphaNumeric(10),    // 10 random alphanumeric characters
		}
	}

	// NOTE : Below is subject for change once we implement the Profile later
	if group == TypeVerification {
		s = []string{
			"M",                              // M for Matchmove
			(config.Instance().AppEnv)[0:2],  // first 2 characters of current app env
			"BO",                             // BO for boilerplate
			group,                            // group
			customRef,                        // custom reference
			strconv.Itoa(time.Now().Year()),  // current year
			strconv.FormatInt(microTime, 10), // timestamp
			GenerateRandAlpha(5),             // 5 random letters
			GenerateUniqueNumeric(5),         // 5 numeric characters
			GenerateRandAlphaNumeric(4),      // 4 random alphanumeric characters
		}
	}

	if group == TypeReferenceSdk {
		s = []string{
			"MBO", // M for Matchmove, BO for boilerplate
			strconv.Itoa(time.Now().Year()) + strconv.FormatInt(microTime, 10), // current year + timestamp
			group + GenerateRandAlphaNumeric(10),                               // group + 5 random letters
			GenerateUniqueNumeric(5),                                           // 5 numeric characters
		}
	}

	if customJoin != "" {
		return strings.Join(s, customJoin)
	}

	return strings.Join(s, "")
}

// TrimStructs trim leading and trailing whitespaces on struct fields
func TrimStructs(v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var structMap map[string]interface{}
	if err := json.Unmarshal(bytes, &structMap); err != nil {
		return err
	}

	structMap = TrimInterface(structMap).(map[string]interface{})
	bytes2, err := json.Marshal(structMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes2, v)
	if err != nil {
		return err
	}

	return nil
}

// TrimInterface trim leading and trailing whitespaces on a given interface (array of interface)
func TrimInterface(data interface{}) interface{} {
	if values, valid := data.([]interface{}); valid {
		for i := range values {
			data.([]interface{})[i] = TrimInterface(values[i])
		}
	} else if values, valid := data.(map[string]interface{}); valid {
		for k, v := range values {
			data.(map[string]interface{})[k] = TrimInterface(v)
		}
	} else if value, valid := data.(string); valid {
		data = strings.TrimSpace(value)
	}

	return data
}

// ConvertUTCTimeToConsumerTimezone from UTC datetime to specified consumer timezone
func ConvertUTCTimeToConsumerTimezone(dateTime time.Time, timezone, layout string) string {
	if timezone == "" {
		timezone = DefaultTimezone
	}

	loc, _ := time.LoadLocation(timezone)
	convertedTime := dateTime.In(loc)

	if layout == "" {
		layout = DefaultDateTimeLayout
	}

	return convertedTime.Format(layout)
}

// CapitalizeFirstChar make first character of a sentence to capital letters
func CapitalizeFirstChar(s string) string {
	for index, value := range s {
		return string(unicode.ToUpper(value)) + s[index+1:]
	}
	return ""
}

// IsMasked checks if required to mask or not
func IsMasked(mask interface{}) bool {
	var isMasking bool
	switch v := mask.(type) {
	case bool:
		isMasking = v
	case string:
		stringVal, err := strconv.ParseBool(v)
		if err != nil {
			isMasking = false
		} else {
			isMasking = stringVal
		}
	case float64:
		intVal := int(v)
		isMasking = intVal != 0
	case int:
		isMasking = v != 0
	default:
		isMasking = false
	}

	return isMasking
}

// Mask mask specific length of a string
func Mask(s string, start, limit int, exceptCharacter []string) string {
	rs := []rune(s)
	for i := start; i < len(rs)-limit; i++ {
		if len(exceptCharacter) > 0 {
			if exists, _ := array.InArray(rs[i], exceptCharacter); exists {
				continue
			}
		}

		rs[i] = 'X'
	}

	return string(rs)
}

// DownloadImage downloads an image from the given URL and returns it as a 2D byte slice ([][]byte)
func DownloadImage(urls []string) ([][]byte, error) {
	imageBytes := make([][]byte, 0)

	for _, url := range urls {
		// Make HTTP request to download the image
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to download image from URL: %w", err)
		}
		defer resp.Body.Close()

		// Check if the response status code is OK
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
		}

		// Read the image data into memory
		imgData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read image data: %w", err)
		}

		// Decode the image (optional if you need to handle image-specific logic)
		_, _, err = image.Decode(bytes.NewReader(imgData))
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %w", err)
		}

		// Convert the image to [][]byte (depending on how you need the image data)
		// Here, I'm assuming you want to convert the entire image byte array into [][]byte.
		imageBytes = append(imageBytes, imgData)
	}

	return imageBytes, nil
}

// GetCacheID return cache identifier
func GetCacheID(keys ...string) string {
	c := config.Instance()
	cacheID := strings.ReplaceAll(c.AppServiceName, "-", "_")
	for _, key := range keys {
		if key != "" {
			cacheID += "_" + key
		}
	}
	return strings.ToLower(cacheID)
}

// generateGUID generates a new GUID and returns it as a string without delimiters
func generateGUID() (string, error) {
	guid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate GUID: %w", err)
	}
	return guid.String(), nil
}

// getFormattedDate returns the current date in yyyy-MM-dd format
func getFormattedDate() string {
	return time.Now().Format("2006-01-02")
}

// generateFormattedGUID generates a processing GUID in the format yyyy-MM-dd + 32-bit GUID without delimiters
func generateFormattedGUID() (string, error) {
	date := getFormattedDate()
	guid, err := generateGUID()
	if err != nil {
		return "", err
	}
	// Remove hyphens from GUID
	guid = guid[0:8] + guid[9:13] + guid[14:18] + guid[19:23] + guid[24:]
	return date + guid, nil
}
func GetSequenceID() (string, error) {
	formattedGUID, err := generateFormattedGUID()
	if err != nil {
		return "", err
	}
	return formattedGUID, nil
}

//

// generateDeterministicCode creates a consistent alphanumeric code from input
func GenerateDeterministicCode(data string, length int) string {
	hash := sha256.Sum256([]byte(data))
	num := new(big.Int).SetBytes(hash[:])

	// Convert the big.Int hash to base62
	base62 := toBase62(num)

	// Pad if necessary
	if len(base62) < length {
		base62 = strings.Repeat("0", length-len(base62)) + base62
	}

	return base62[:length]
}

// toBase62 encodes a big.Int to a base62 string
func toBase62(num *big.Int) string {
	if num.Cmp(big.NewInt(0)) == 0 {
		return string(charset[0])
	}

	var result []byte
	base := big.NewInt(int64(len(charset)))
	zero := big.NewInt(0)
	mod := new(big.Int)

	for num.Cmp(zero) > 0 {
		num, mod = new(big.Int).DivMod(num, base, mod)
		result = append([]byte{charset[mod.Int64()]}, result...)
	}

	return string(result)
}

func RemoveMargineFromProductCostPrice(productId int64, cost float64) (costPrice float32) {
	for _, margine := range constant.MARGINEDETAIL {
		if margine.ProductId == productId {
			switch margine.MargineType {
			case "FLATVALUE":
				costPrice = float32(cost - float64(margine.Value))
			case "PERCENTAGE":
				mergineval := float64(1 + margine.Value/100)
				costPrice = float32(math.Floor(cost / mergineval))
			}
			break
		}
	}

	return costPrice
}
