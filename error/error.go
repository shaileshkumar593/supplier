package error

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	customContext "swallow-supplier/context"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils"
)

// ResponseError represents the error response
type ResponseError struct {
	ID        string      `json:"id"`
	Status    int         `json:"status"`
	Code      string      `json:"code"`
	Title     string      `json:"title"`
	Detail    string      `json:"detail"`
	Source    interface{} `json:"source,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// NewError this formats the required error response
func NewError(ctx context.Context, code string, message string, source interface{}) *ResponseError {
	e, _ := GetErrorByCode(code)

	if message != "" {
		message = utils.CapitalizeFirstChar(message)

		// separate by .
		if strings.Contains(message, ".") {
			var sMessage = strings.Split(message, ". ")
			var finalMessage []string
			for _, v := range sMessage {
				vMessage := utils.CapitalizeFirstChar(v)
				finalMessage = append(finalMessage, vMessage)
			}

			message = strings.Join(finalMessage, ". ")
		}
	}
	return &ResponseError{
		ID:        customContext.GetCtxHeader(ctx, customContext.CtxLabelRequestID),
		Status:    e.status,
		Code:      e.code,
		Title:     e.message,
		Detail:    message,
		Source:    source,
		Timestamp: time.Now().UTC().String(),
	}
}

// NewErrorCustom can be used custom error where error details are not specified on the errors.go
func NewErrorCustom(ctx context.Context, code string, message string, title string, statusCode int, source interface{}) *ResponseError {
	return &ResponseError{
		ID:        customContext.GetCtxHeader(ctx, customContext.CtxLabelRequestID),
		Status:    statusCode,
		Code:      code,
		Title:     title,
		Detail:    message,
		Source:    source,
		Timestamp: time.Now().UTC().String(),
	}
}

func (e *ResponseError) Error() string {
	return e.Code
}

// Response return when error raised in the code
func EncodeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	if w == nil {
		log.Println("Error: Response writer is nil")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var (
		code     = http.StatusInternalServerError // Default to internal server error
		response *ResponseError                   // Placeholder for the response
	)

	// Attempt to cast the error to *ResponseError
	if respErr, ok := err.(*ResponseError); ok && respErr != nil {
		if respErr.Status == 0 {
			log.Println("Error: Status code is empty")
			response = NewError(ctx, "invalid_status", "Error status code is empty", nil)
		} else {
			response = respErr
			code = respErr.Status
		}
	} else {
		response = NewError(ctx, "generic_error", "An unexpected error occurred", nil)
	}

	w.WriteHeader(code)
	errorResponse := map[string]interface{}{
		"trace_id": utils.GenerateUUID("", true),
		"body": yanolja.ErrorMsg{
			ErrorCode:    fmt.Sprintf("%d", code),
			ErrorMessage: response.Detail,
		},
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// / EncodeErrorResponseWithBody encodes the error response along with a body
func EncodeErrorResponseWithBody(ctx context.Context, err error, w http.ResponseWriter, responseBody interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var (
		code     = http.StatusInternalServerError
		response *ResponseError
	)

	if respErr, ok := err.(*ResponseError); ok && respErr != nil {
		response = respErr
		code = response.Status
	} else {
		response = NewError(ctx, "generic_error", "An unexpected error occurred", nil)
	}

	w.WriteHeader(code)

	/* errorResponse := map[string]interface{}{
		"status": code,
		"code":   response.Code,
		"detail": response.Detail,
		"body":   responseBody,
	} */

	// when using  this variable test  code works fine
	errorResponse := map[string]interface{}{
		"trace_id": utils.GenerateUUID("", true),
		"body": yanolja.ErrorMsg{
			ErrorCode:    fmt.Sprintf("%d", code),
			ErrorMessage: fmt.Sprint(response.Detail + "  " + responseBody.(string)),
		},
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
