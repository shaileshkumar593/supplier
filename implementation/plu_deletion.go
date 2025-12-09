package implementation

import (
	"context"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// DeleteExpiredPlu based on
// 1. based on invalid version, invalid date, invalid product, productvariant
func (s *service) DeleteExpiredPlu(ctx context.Context, req yanolja.AllProduct) (resp common.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "GetProducts",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	return
}
