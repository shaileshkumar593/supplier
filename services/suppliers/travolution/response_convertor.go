package travolution

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func ResponseConvertor(ctx context.Context, response *client.Response, logger log.Logger, err1 error) (res travolution.Response, err error) {

	if response == nil {
		res.Body = "nil response"
		res.Code = "500"
		return res, customError.NewError(ctx, "leisure-api-0001", "nil response", nil)
	}
	res.Code = strconv.Itoa(response.Status)

	if err == nil && response.Status == 403 {
		res.Code = strconv.Itoa(http.StatusForbidden)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), ServiceName), nil)
		return res, err
	}

	if err == nil && response.Status == 503 {
		res.Code = strconv.Itoa(http.StatusServiceUnavailable)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-0004", fmt.Sprintf(customError.ErrExternalService.Error(), ServiceName), nil)
		return res, err
	}

	if err == nil && response == nil && response.Status != http.StatusOK {
		level.Info(logger).Log("err", response)
		return res, customError.NewError(ctx, res.Code, "empty response", nil)
	}

	level.Info(logger).Log("info", "response body ", response.Body)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(ctx, res.Code, "Empty Body", nil)
	}

	if unmarshalErr := json.Unmarshal([]byte(response.Body), &res.Body); unmarshalErr != nil {
		level.Info(logger).Log("response unmarshall err", unmarshalErr)

		return res, customError.NewError(ctx, res.Code, "Invalid JSON Body", unmarshalErr)
	}

	return res, nil
}
