package implementation

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/yanolja"

	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// not needed now because we are doing simple image url sync and getting ImageId
// InventrySync for syncing inventory to trip
func (s *service) UpdateImageSyncStatus(ctx context.Context, req []domain.ImageUrlForProcessing) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "InventrySync",
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

	err = s.mongoRepository[config.Instance().MongoDBName].BulkUpdateProductImageStatusAndImageId(ctx, req)
	if err != nil {
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}
	// write code to add imagewithId  to mongo database as seprate collection
	resp.Code = "200"
	resp.Body = fmt.Sprintln("status updated successfully")
	return resp, nil
}
