package v2_test

import (
	"os"
	"testing"

	validator "swallow-supplier/utils/validator/v2"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	validator.Init()
	os.Exit(m.Run())
}

func TestValidatorWithSuccess(t *testing.T) {
	errorMessage := `required field`
	type test struct {
		Name string `validate:"required"`
	}

	req := test{}

	res, err := validator.Struct(req)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, errorMessage, res[0].(map[string][]interface{})[`Name`][0].(string))
}

func TestValidatorInvalidTag(t *testing.T) {
	type test struct {
		Name string `validate:"require"`
	}

	req := test{}

	res, err := validator.Struct(req)
	assert.Error(t, err)
	assert.Empty(t, res)
	assert.Equal(t, `Undefined validation function 'require' on field 'Name'`, err.Error())
}

func TestValidatorNilRequest(t *testing.T) {
	res, err := validator.Struct(nil)
	assert.Error(t, err)
	assert.Empty(t, res)
}

func TestValidatorWithFormattedErrors(t *testing.T) {
	errorMessage := `required field`
	type test struct {
		Name   string `json:"name" validate:"required"`
		Mobile string `json:"mobile" validate:"len=10"`
		Alpha  string `json:"alpha" validate:"alphanumchar"`
		// json as -
		Address string `json:"-" validate:"required"`
		// unregistered error tag
		Hostname string `json:"hostname" validate:"hostname"`
	}

	req := test{}

	res, err := validator.StructWithFormattedErrors(req)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, errorMessage, res[`name`][0].(string))
	assert.Equal(t, `hostname`, res[`hostname`][0].(string))
	assert.Equal(t, `required length of 10`, res[`mobile`][0].(string))
	assert.Equal(t, errorMessage, res[`Address`][0].(string))
	assert.Equal(t, `invalid string format`, res[`alpha`][0].(string))
}

func TestValidatorWithFormattedErrorsNilInput(t *testing.T) {
	res, err := validator.StructWithFormattedErrors(nil)
	assert.Error(t, err)
	assert.Empty(t, res)
}
