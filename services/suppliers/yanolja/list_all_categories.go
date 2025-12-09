package yanolja

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetAllCategories get the  all categories for T&A
func (y *Yanolja) GetAllCategories(ctx context.Context) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "GetAllCategories")
	level.Info(logger).Log("info", "Service GetAllCategories")

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+"/v1/categories",
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)

	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		if response.Status == http.StatusServiceUnavailable {
			return res, customError.NewError(y.Ctx, "external_processing_error", " Service Temporarily Unavailable", nil)
		}
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "empty response", nil)
	}

	level.Info(logger).Log("info", "response body ", response.Body)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		level.Error(logger).Log("error", "unmarshal error")
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}
	return res, nil
}
