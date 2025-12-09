package implementation

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils"
)

// HeartBeat endpoint
// Checks the availability of the Service
func (s *service) HeartBeat(ctx context.Context) (res yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)
	res.TraceID = utils.GenerateUUID("", true)

	logger := log.With(
		s.logger,
		"method", "heartbeat",
		"Request ID", requestID,
		"Trace ID", res.TraceID,
	)

	resrepo, err := s.mongoRepository[cf.MongoDBName].GetHeartBeatFromMongo(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
		res.Code = "503"
		res.Body = resrepo
		return res, customError.NewError(ctx, "leisure-api-1015", err.Error(), nil)
	}
	level.Info(logger).Log("res", res)
	res.Body = resrepo
	res.Code = "200"
	return res, nil
}
