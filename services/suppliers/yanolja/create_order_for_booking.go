package yanolja

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// WaitingForOrder created order in waiting state
func (y *Yanolja) WaitingForOrder(ctx context.Context, req yanolja.WaitingForOrder) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "WaitingForOrder")
	level.Info(logger).Log("info", "Service WaitingForOrder")
	host := y.Host + "/v1/orders/prepare"
	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+"/v1/orders/prepare",
		http.MethodPost,
		client.ContentTypeJSON,
		req,
	)
	level.Error(logger).Log("info", "url", host, "response", response, "err", err)

	return YanoljaResponseConversion(y.Ctx, response, logger, err)
}

// OrderComplete WaitingForOrder created order in waiting state
func (y *Yanolja) OrderComplete(ctx context.Context, req yanolja.OrderConfirmation) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "OrderComplete")
	level.Info(logger).Log("info", "Calling Yanolja /v1/orders/complete")
	host := y.Host + "/v1/orders/complete"
	//time.Sleep(310 * time.Second)
	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+"/v1/orders/complete",
		http.MethodPost,
		client.ContentTypeJSON,
		req,
	)
	level.Error(logger).Log("info", "***OrderComplete***XXXXXX----complete******", "url", host, response, "err", err)
	// fmt.Println("///////////// response //////////////// ", response)
	// fmt.Println("########### response body ############# ", string(response.Body))
	// fmt.Println("<<<<<<<<<<<< error >>>>>>>>>>>>>>>> ", err)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			level.Error(logger).Log("error", "context canceled during HTTP call", "ctxErr", ctx.Err())
			res, err = y.TimeoutAndForcedOrderCancellation(ctx, req.PartnerOrderId)
			return res, err

		} else if errors.Is(err, context.DeadlineExceeded) {
			logger.Log("msg", "HTTP context deadline exceeded", "error", err.Error())
			res, err = y.TimeoutAndForcedOrderCancellation(ctx, req.PartnerOrderId)
			return res, err
		}

		level.Error(logger).Log("error", "failed Yanolja call", "err", err)
	}

	if err == nil && response.Status == 504 && strings.Contains(string(response.Body), "Gateway Time-out") {
		level.Info(logger).Log(".....cancelllation time out request called ......")
		_, _ = y.TimeoutAndForcedOrderCancellation(ctx, req.PartnerOrderId)
		res.Code = strconv.Itoa(http.StatusServiceUnavailable)
		res.Body = response.Body
		err = customError.NewError(ctx, "leisure-api-00025", fmt.Sprintf(customError.ErrGatewayTimeout.Error(), ServiceName), nil)
		return res, err
	}

	if err == nil && response.Body != "" && response.Status != http.StatusOK {
		level.Info(logger).Log("info", " timeout in ggt checket ")
		// Declare a variable to hold the map
		result := make(map[string]interface{})

		// Unmarshal the JSON string into the map
		err = json.Unmarshal([]byte(response.Body), &result)
		if err != nil {
		}
		// Check if "body" exists in the map
		body, _ := result["body"].(map[string]interface{})

		// Access the "message" key in the "body" map
		code, _ := body["code"].(string)

		if code == "leisure-api-0012" {
			_, _ = y.ForcedOrderCancellation(ctx, req.PartnerOrderId)
		}

		level.Error(logger).Log("error ", response.Body)
		text, _ := ErrorTextConversion(ctx, response, logger)
		res.Body = text
		return res, customError.NewError(ctx, code, text, ServiceName)
	}

	return YanoljaResponseConversion(y.Ctx, response, logger, err)
}

// GetOrderbyId get order based on the orderId
func (y *Yanolja) GetOrderbyId(ctx context.Context, req yanolja.OrderConfirmation) (res yanolja.Response, err error) {

	logger := log.With(y.Service.Logger, "method", "GetOrderbyId")
	level.Info(logger).Log("info", "Service GetOrderbyId")

	urls := fmt.Sprintf(`/v1/orders/%d?partnerOrderId=%s`, req.OrderId, req.PartnerOrderId)
	level.Info(logger).Log("info", "url for get request ", "urls :", urls)

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	level.Error(logger).Log("info", "********GetOrderById************* ", "url", urls, "response : ", response, "err", err)
	return YanoljaResponseConversion(y.Ctx, response, logger, err)
}

// CancelFullOrder cancel order based on the orderId
func (y *Yanolja) CancelFullOrder(ctx context.Context, orderid int64) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "CancelFullOrder")
	level.Info(logger).Log("info", "Service CancelFullOrder")

	urls := fmt.Sprintf(`/v1/orders/%d/full-cancel`, orderid)
	level.Info(logger).Log("info", "url for get request ", urls)

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodPost,
		client.ContentTypeJSON,
		nil,
	)

	level.Error(logger).Log("info", "**CancelFullOrder***XXXXXX----fullCancel***** ", "url", urls, "response", response, "err", err)

	return YanoljaResponseConversion(y.Ctx, response, logger, err)
}

// TimeoutAndForcedOrderCancellation use for forced cancellation
func (y *Yanolja) TimeoutAndForcedOrderCancellation(ctx context.Context, partnerOrderId string) (res yanolja.Response, err error) {

	logger := log.With(y.Service.Logger, "method", "TimeoutAndForcedOrderCancellation")
	level.Info(logger).Log("info", "Service TimeoutAndForcedOrderCancellation")

	urls := fmt.Sprintf(`/v1/orders/timeout-cancel?partnerOrderId=%s`, partnerOrderId)
	level.Info(logger).Log("info", "url for get request ", urls)

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodPost,
		client.ContentTypeJSON,
		nil,
	)
	level.Error(logger).Log("info", "********TimeoutAndForcedOrderCancellation*****XXXXXX----TimeoutForcedCancel********* ", "url", urls, "response", response, "err", err)
	return YanoljaResponseConversion(y.Ctx, response, logger, err)
}

// ForcedOrderCancellation use for forced cancellation
func (y *Yanolja) ForcedOrderCancellation(ctx context.Context, partnerOrderId string) (res yanolja.Response, err error) {

	logger := log.With(y.Service.Logger, "method", "ForcedOrderCancellation")
	level.Info(logger).Log("info", "Service ForcedOrderCancellation")

	urls := fmt.Sprintf(`/v1/orders/force-cancel?partnerOrderId=%s`, partnerOrderId)
	level.Info(logger).Log("info", "url for get request ", urls)

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodPost,
		client.ContentTypeJSON,
		nil,
	)
	level.Error(logger).Log("info", "********ForcedOrderCancellation*****XXXXXX----ForcedCancel********* ", "url", urls, "response", response, "err", err)

	return YanoljaResponseConversion(y.Ctx, response, logger, err)
}
