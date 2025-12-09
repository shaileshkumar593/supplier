package validator

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/validator.v2"

	"swallow-supplier/utils/array"
)

const (
	// PatternEmail regexp pattern for RFC 5322 (email) electronic mail address
	PatternEmail = `(?:[a-z0-9!#$%&'*+/=?^_` + "`" +
		`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" +
		`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`

	// PatternString a-z, A-Z, 0-9, /-?:().,'+.#!@&
	PatternString = `[#.0-9a-zA-Z ,\\/\\_:+?')(@#!&-]`
)

// Init the validator
func Init() {
	validator.SetValidationFunc("required", Required)
	validator.SetValidationFunc("maxlen", MaxLen)
	validator.SetValidationFunc("minlen", MinLen)
	validator.SetValidationFunc("reqlen", ReqLen)
	validator.SetValidationFunc("enum", Enum)
	validator.SetValidationFunc("notempty", NotEmpty)
	validator.SetValidationFunc("date", Date)
	validator.SetValidationFunc("url", URL)
	validator.SetValidationFunc("numeric", Numeric)
	validator.SetValidationFunc("pattern", Pattern)
	validator.SetValidationFunc("email", Email)
	validator.SetValidationFunc("alphanumchar", AlphaNumChar)
}

// Required validate if field exist
func Required(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return errors.New("required field")
	}

	return nil
}

// MaxLen - Max Length validation
func MaxLen(v interface{}, param string) (err error) {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return nil
	}
	sl, _ := strconv.Atoi(param)
	if sl < len(st.String()) {
		return errors.New("maximum length allowed is " + param)
	}

	return nil
}

// MinLen - Min Length validation
func MinLen(v interface{}, param string) (err error) {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return nil
	}
	sl, _ := strconv.Atoi(param)
	if sl > len(st.String()) {
		return errors.New("minimum length allowed is " + param)
	}

	return nil
}

// ReqLen - Required Length validation
func ReqLen(v interface{}, param string) (err error) {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return nil
	}
	sl, _ := strconv.Atoi(param)
	if sl != len(st.String()) {
		return errors.New("required length of " + param)
	}

	return nil
}

// Enum - Enumeration validation
func Enum(v interface{}, param string) (err error) {
	st := reflect.ValueOf(v).String()
	if st == "" {
		return nil
	}
	options := strings.Split(param, "|")
	if exist, _ := array.InArray(st, options); !exist {
		return fmt.Errorf("values supported: [%s]", strings.Join(options, "|"))
	}
	return nil
}

// NotEmpty - checks if slice has value
func NotEmpty(v interface{}, param string) (err error) {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		if st.Len() == 0 {
			sliceType := reflect.TypeOf(v).Elem()
			return errors.New("required to have 1 or more " + strings.ToLower(sliceType.Name()) + " values")
		}
	}
	return nil
}

// Date - checks if valid date format
func Date(v interface{}, param string) (err error) {
	val := reflect.ValueOf(v).String()
	if val == "" {
		return nil
	}

	var l string

	switch param {
	case "2006-01-02":
		l = "YYYY-MM-DD"
		break
	}

	_, err = time.Parse(param, val)
	if err != nil {
		return errors.New("expected format of " + l)
	}

	return nil
}

// URL - checks if a valid url
func URL(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return nil
	}

	u, err := url.ParseRequestURI(st.String())
	if err != nil {
		return errors.New("should be a valid url")
	}

	if strings.HasPrefix(u.Host, ".") || err != nil {
		return errors.New("should be a valid url")
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return errors.New("should be a valid url")
	}
	if len(u.Scheme) == 0 {
		return errors.New("should be a valid url")
	}

	return nil
}

// Numeric - Number only validation
func Numeric(v interface{}, param string) (err error) {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return nil
	}

	_, err = strconv.ParseFloat(st.String(), 64)
	if err != nil {
		return errors.New("should be a valid numeric value")
	}

	return nil
}

// Pattern regex validation
func Pattern(v interface{}, param string) (err error) {
	val := reflect.ValueOf(v).String()
	if val == "" {
		return nil
	}

	if ok, _ := regexp.MatchString(`^(`+param+`)$`, val); !ok {
		return errors.New("value should be on format : " + param)
	}

	return nil
}

// Email format validation
func Email(v interface{}, param string) (err error) {
	val := reflect.ValueOf(v).String()
	if val == "" {
		return nil
	}

	if ok, _ := regexp.MatchString(`^(`+PatternEmail+`)$`, val); !ok {
		return errors.New("invalid email format")
	}

	return nil
}

// AlphaNumChar format validation
func AlphaNumChar(v interface{}, param string) (err error) {
	val := reflect.ValueOf(v).String()
	if val == "" {
		return nil
	}

	var allSupported = regexp.MustCompile(`^` + PatternString + `+$`).MatchString(val)
	if !allSupported {
		return errors.New("invalid string format")
	}

	return nil
}
