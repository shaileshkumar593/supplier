package implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	model "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/yanolja"
	yanoljasvc "swallow-supplier/services/suppliers/yanolja"
	"swallow-supplier/utils"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// GetRegions get all products from yanolja
func (s *service) GetRegions(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetRegions",
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

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetRegionalCategories(ctx)

	if err != nil {
		e := level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		return resp, e
	}

	level.Info(logger).Log("response", resp)

	return resp, nil
}

// InsertAllRegions insert all regions from yanolja
func (s *service) InsertAllRegions(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "InsertAllRegions",
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

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetRegionalCategories(ctx)

	if err != nil {
		e := level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		return resp, e
	}

	level.Info(logger).Log("response", resp)
	doc, err := json.Marshal(resp.Body)
	if err != nil {
		resp.Code = "500"
		return resp, fmt.Errorf("marshal error ")
	}

	// Unmarshal JSON response
	var rec []model.Region
	err = json.Unmarshal([]byte(doc), &rec)
	if err != nil {
		level.Error(logger).Log("error", "failed to unmarshal JSON response", err)
		resp.Code = "500"
		return resp, fmt.Errorf("json unmarshal error: %w", err)
	}

	err = s.mongoRepository[config.Instance().MongoDBName].InsertRegions(ctx, rec)
	if err != nil {
		level.Error(logger).Log("error", "request to InsertRegions raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from InsertRegions, %v", err), "InsertRegions")
	}

	return resp, nil

}
