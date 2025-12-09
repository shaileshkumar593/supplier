package trip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	customError "swallow-supplier/error"
	"swallow-supplier/mongo/domain/trip"
	"swallow-supplier/request_response/yanolja"

	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetOrderbyId get order based on the orderId
func (tp *Trip) NotifyToTrip(ctx context.Context, payload trip.CallBackRequest) (res yanolja.Response, err error) {

	logger := log.With(tp.Service.Logger, "method", "NotifyToTrip")

	level.Info(logger).Log("payload ", payload)
	level.Info(logger).Log("endpoint of trip :", tp.Host+"/out/notify")

	response, err := tp.Service.Send(
		ctx,
		ServiceName,
		tp.Host+"/out/notify",
		http.MethodPost,
		client.ContentTypeJSON,
		payload,
	)
	if err != nil {
		level.Error(logger).Log("error ", "trip service error")
		res.Body = response.Body
		res.Code = string(response.Status)
		return res, customError.NewError(ctx, "leisure-api-1022", fmt.Sprintf("trip request error %v", err), "NotifyToTrip")
	}

	return TripResponseConversion(tp.Ctx, response, logger, err)
}

func TripResponseConversion(ctx context.Context, response *client.Response, logger log.Logger, err1 error) (res yanolja.Response, err error) {
	res.Code = strconv.Itoa(response.Status)

	if response.Status == 403 {
		res.Code = strconv.Itoa(http.StatusForbidden)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), ServiceName), nil)
		return res, err
	}

	if err1 == nil && response.Status == 503 {
		res.Code = strconv.Itoa(http.StatusServiceUnavailable)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-0004", fmt.Sprintf(customError.ErrExternalService.Error(), ServiceName), nil)
		return res, err
	}

	if err1 == nil && response.Body == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Error(logger).Log("err", string(response))
		return res, customError.NewError(ctx, res.Code, "empty response", nil)
	}

	if err1 == nil && response.Body != "" && response.Status != http.StatusOK {
		level.Error(logger).Log("error ", response.Body)
		res.Body = fmt.Sprintln(response.Message, response.Body)
		res.Code = string(response.Status)
		return res, customError.NewError(ctx, "leisure-api-0012", res.Body.(string), nil)
	}

	level.Info(logger).Log("info", "response body ", response.Body, err1)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(ctx, res.Code, "Empty Body", nil)
	}

	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		return res, customError.NewError(ctx, res.Code, "Empty Body", nil)
	}

	level.Info(logger).Log(" trip response ", res)

	if res.Code != "200" {
		jsonstr, err2 := json.Marshal(res.Body)
		if err2 != nil {
			level.Error(logger).Log("error", "Error in marshaling data:", err2)
			res.Body = "Error in marshaling data"
			res.Code = "500"
			return res, err2
		}
		var resBody yanolja.ResponseBody
		// Unmarshal the JSON data into the struct
		err2 = json.Unmarshal(jsonstr, &resBody)
		if err2 != nil {
			level.Error(logger).Log("error", "Error in unmarshaling data:", err2)
			res.Body = "Error in unmarshaling data :"
			res.Code = "500"
			return res, err2
		}

		return res, customError.NewErrorCustom(ctx, res.Code, resBody.Detail, resBody.Message, response.Status, ServiceName)
	}

	return res, nil
}
