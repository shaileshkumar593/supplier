package implementation

import (
	"context"
	"fmt"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/travolution"
	travolutionSvc "swallow-supplier/services/suppliers/travolution"
	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func (s *service) GetAdditionalInfos(ctx context.Context, req travolution.BookingAdditionalInfoRequest) (resp travolution.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetAdditionalInfos",
		"Request ID", requestID,
	)
	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	level.Info(logger).Log(" info ", "travolution service call")

	var tsvc, _ = travolutionSvc.New(ctx)
	resp, err = tsvc.GetBookingAdditionalInfo(ctx, req)
	if err != nil {
		level.Error(logger).Log("treavolution error ", err)

		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprint(customError.ErrForbiddenClient.Error(), "GetProductByUid"), nil)
		} else {
			err = fmt.Errorf("request to travolution client raised error")
			resp.Code = "500"
			resp.Body = err
		}
		return resp, err
	}

	level.Info(logger).Log("response ", resp)
	return resp, nil
}

// GetAdditionalInfoByUid
func (s *service) GetAdditionalInfoByUid(ctx context.Context, req travolution.BookingAdditionalInfoRequest) (resp travolution.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetAdditionalInfoByUid",
		"Request ID", requestID,
	)
	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	level.Info(logger).Log(" info ", "travolution service call")

	var tsvc, _ = travolutionSvc.New(ctx)
	resp, err = tsvc.GetBookingAdditionalInfo(ctx, req)
	if err != nil {
		level.Error(logger).Log("treavolution error ", err)

		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprint(customError.ErrForbiddenClient.Error(), "GetProductByUid"), nil)
		} else {
			err = fmt.Errorf("request to travolution client raised error")
			resp.Code = "500"
			resp.Body = err
		}
		return resp, err
	}
	resp.Code = "200"

	level.Info(logger).Log("response ", resp)
	return resp, nil
}
