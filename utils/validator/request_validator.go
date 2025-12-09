package validator

import (
	"context"
	"fmt"

	customError "swallow-supplier/error"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
)

// add switch case based on type of variable passed

// ValidateRequest validates the incoming request and logs validation errors if present.
func ValidateRequest(ctx context.Context, logger log.Logger, req interface{}) (err error) {

	level.Info(logger).Log("info", "request validator ")
	validate := validator.New()

	// Validate the request struct
	err = validate.Struct(req)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			source := make(map[string][]interface{})
			compiledErrors := make([]interface{}, 0)

			// Iterate over each validation error
			for _, ve := range validationErrors {
				field := ve.Field()                                                     // Get the struct field name
				errMsg := fmt.Sprintf("Error: %s, Condition: %s", ve.Tag(), ve.Param()) // Example error message
				source[field] = append(source[field], errMsg)
			}

			compiledErrors = append(compiledErrors, source)

			// Return formatted validation error response
			return customError.NewError(ctx, "invalid_model", customError.ErrValidation.Error(), FormatCompiledErrors(compiledErrors))
		}

		// Handle unexpected validation errors
		return customError.NewError(ctx, "validation_error", err.Error(), nil)
	}

	return nil
}

// FormatCompiledErrors formats error response
func FormatCompiledErrors(errs []interface{}) map[string][]interface{} {
	var errSource = make(map[string][]interface{})
	for _, e := range errs {
		if e.(map[string][]interface{}) != nil {
			for f, ex := range e.(map[string][]interface{}) {
				errSource[f] = ex
			}
		}
	}

	return errSource
}
