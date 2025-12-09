package v2

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	validate "github.com/go-playground/validator/v10"
)

var (
	val                *validate.Validate
	customErrorMessage = map[string]string{
		`required`:     `required field`,
		`len`:          `required length of %v`,
		`min`:          `minimum length allowed is %v`,
		`max`:          `maximum length allowed is %v`,
		`oneof`:        `values supported: [%v]`, // enum either one of them
		`alphanumchar`: `invalid string format`,
	}
)

const (
	// PatternString a-z, A-Z, 0-9, /-?:().,'+.#!@&
	PatternString = `[#.0-9a-zA-Z ,\\/\\_:+?')(@#!&-]`

	// AlphaNumCharTag is custom validaton tag
	AlphaNumCharTag = `alphanumchar`
)

// Init the validator
func Init() {
	val = validate.New()
	val.RegisterValidation(AlphaNumCharTag, AlphaNumChar)
	// This function will be used to get the json tag from the struct
	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// New will intialize and return validate
func New() *validate.Validate {
	if val != nil {
		return val
	}
	Init()
	return val
}

// Struct will return compiled error list
func Struct(req interface{}) (compiledErrors []interface{}, err error) {
	var (
		source = make(map[string][]interface{})
	)
	defer func() {
		if errInternal := recover(); errInternal != nil {
			errd := errInternal.(string)
			err = fmt.Errorf(`%v`, errd)
			return
		}
	}()
	v := New()
	err = v.Struct(req)
	if err != nil {
		// Below error check is for nil interface
		if _, ok := err.(*validate.InvalidValidationError); ok {
			return nil, err
		}

		for _, err := range err.(validate.ValidationErrors) {
			// this steps is used for removing struct name from the namespace
			arr := strings.Split(err.Namespace(), `.`)
			f := strings.Join(arr[1:], `.`)
			source[f] = append(source[f], getCustomErrorMessage(err.Tag(), err.Param()))
		}
		compiledErrors = append(compiledErrors, source)
	}
	return compiledErrors, nil
}

// StructWithFormattedErrors will return formatted errors
func StructWithFormattedErrors(req interface{}) (map[string][]interface{}, error) {
	errs, err := Struct(req)
	if err != nil {
		return nil, err
	}
	var compiledErrors = make(map[string][]interface{})
	for _, e := range errs {
		if e.(map[string][]interface{}) != nil {
			for f, ex := range e.(map[string][]interface{}) {
				compiledErrors[f] = ex
			}
		}
	}
	return compiledErrors, nil
}

func AlphaNumChar(fl validate.FieldLevel) bool {
	return regexp.MustCompile(`^` + PatternString + `+$`).MatchString(fl.Field().String())
}

func getCustomErrorMessage(tagName, param string) string {
	if customErrorMessage[tagName] == "" {
		return tagName
	}

	if param != "" {
		return fmt.Sprintf(customErrorMessage[tagName], param)
	}
	return customErrorMessage[tagName]
}
