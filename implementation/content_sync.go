package implementation

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"

	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ProductContentSync  call to sync product data to trip
func (s *service) ProductContentSync(ctx context.Context) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "ProductContentSync",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Error(logger).Log("error", "processing request went into panic mode", "panic", r)
		resp.Code = "500"
		err = fmt.Errorf("panic occurred: %v", r)

	}(ctx)

	productcontent, err := s.mongoRepository[config.Instance().MongoDBName].GetProductContentNotSync(ctx)
	if err != nil {
		level.Error(logger).Log("error", " repository error on request to GetProductContentNotSync", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update by orderId, %v", err), "GetProductContentNotSync")
	}

	contentReq := trip.ProductContentSync{
		Message: "ProductContent",
		Data:    productcontent,
	}
	resp.Code = "200"
	resp.Body = contentReq

	return resp, nil
}

// PackageContentSync  call to sync package data to trip
func (s *service) PackageContentSync(ctx context.Context) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PackageContentSync",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Error(logger).Log("error", "processing request went into panic mode", "panic", r)
		resp.Code = "500"
		err = fmt.Errorf("panic occurred: %v", r)

	}(ctx)

	packagecontent, err := s.mongoRepository[config.Instance().MongoDBName].GetPackageContentNotSync(ctx)
	if err != nil {
		level.Error(logger).Log("error", " repository error on request to GetPackageContentNotSync", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update by orderId, %v", err), "GetPackageContentNotSync")
	}

	contentReq := trip.PackageContentSync{
		Message: "PackagetContent",
		Data:    packagecontent,
	}
	resp.Code = "200"
	resp.Body = contentReq

	return resp, nil

}

// UpdateContentSyncStatus
func (s *service) UpdateContentSyncStatus(ctx context.Context, contentSyncStatus trip.TripMessageForSync) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "UpdateContentSyncStatus",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Error(logger).Log("error", "processing request went into panic mode", "panic", r)
		resp.Code = "500"
		err = fmt.Errorf("panic occurred: %v", r)

	}(ctx)

	err = s.mongoRepository[config.Instance().MongoDBName].BulkUpdateSyncStatus(ctx, contentSyncStatus)
	if err != nil {
		level.Error(logger).Log("error", " repository error on request to BulkUpdateSyncStatus", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update by orderId, %v", err), "BulkUpdateSyncStatus")
	}

	resp.Code = "200"
	resp.Body = fmt.Sprintln("package sync status updated")

	return resp, nil
}
