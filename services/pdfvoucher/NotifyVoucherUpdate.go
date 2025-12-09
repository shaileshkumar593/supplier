package pdfvoucher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetOrderbyId get order based on the orderId
func (vp *VoucherPdf) NotifyToVoucherPdfUpdate(ctx context.Context, payload []common.PdfVoucherRequest) (res yanolja.Response, err error) {

	logger := log.With(vp.Service.Logger, "method", "NotifyToVoucherPdfUpdate")

	level.Info(logger).Log("payload ", payload)
	level.Info(logger).Log("endpoint of pdf voucher :", vp.Host+"/create/pdf")

	response, err := vp.Service.Send(
		ctx,
		ServiceName,
		vp.Host+"/create/pdf",
		http.MethodPost,
		client.ContentTypeJSON,
		payload,
	)
	if err != nil {
		level.Error(logger).Log("error ", "pdf voucher service error")
		res.Body = response.Body
		res.Code = string(response.Status)
		return res, customError.NewError(ctx, "leisure-api-1022", fmt.Sprintf("pdf voucher request error %v", err), "NotifyToTrip")
	}

	fmt.Println("response body   ", string(response.Body))

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
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(ctx, res.Code, "empty response", nil)
	}

	level.Info(logger).Log("info", "response body ", response.Body)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(ctx, res.Code, "Empty Body", nil)
	}

	res.Body = response.Body
	level.Info(logger).Log("response from pdf generator ", res)
	return res, nil
}
