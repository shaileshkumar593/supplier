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

// GetRegionalCategories get the list of regional category
func (y *Yanolja) GetRegionalCategories(ctx context.Context) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "GetRegionalCategories")
	level.Info(logger).Log("info", "Service GetProduct")

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+"/v1/regions",
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)

	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "empty response", nil)
	}

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
		level.Error(logger).Log("error", "unmarshal error")
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	return res, nil
}
