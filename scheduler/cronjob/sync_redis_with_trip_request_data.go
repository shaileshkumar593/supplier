package cronjob

import (
	"context"
	"fmt"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	svc "swallow-supplier/iface"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// SyncTripRequestToredis  to sync sequenceId to redis
func SyncTripRequestToredis(ctx context.Context, mrepo svc.MongoRepository, logger log.Logger) (err error) {
	/* level.Info(logger).Log(
		"method name ", "SchdulePluUpsertToRedis",
	) */

	trip_Request, err := mrepo.FetchTripRequests(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "Error in FetchTripRequests", "error", err)
		return fmt.Errorf("error in FetchTripRequests: %w", err)
	}
	if len(trip_Request) == 0 {
		level.Error(logger).Log("msg", "no order to sync to redis", "error", err)
		return fmt.Errorf("no order to sync to redis")
	}
	//fmt.Println("**************** FetchTripRequests  is :******************", trip_Request)

	// check for empty mongo repo data
	var allEmpty bool = true
	for _, reqVal := range trip_Request {

		if len(reqVal) > 0 {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		level.Info(logger).Log("Info  ", " no order is created by trip till now")
		return fmt.Errorf("no order is created by trip till now")
	}

	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("msg", "Error initializing cache layer", "error", err)
		return fmt.Errorf("error initializing cache layer: %w", err)
	}

	for _, reqdata := range trip_Request {

		for _, data := range reqdata {
			key := data.OtaOrderID + "-" + data.RequestCategory
			val := data.SequenceID
			exist := cacheLayer.Exist(ctx, key)
			if exist == 0 {
				fmt.Println("key : ", key, "  value : ", val)
				// Insert or update the record in Redis
				if err := cacheLayer.Set(ctx, key, val); err != nil {
					level.Error(logger).Log(
						"msg", "Failed to set key in Redis for every second",
						"key", key,
						"value", val, // Log the serialized JSON for debugging
						"error", err,
					)
					return fmt.Errorf("failed to set key %s in Redis: %w", key, err)
				}
			}

		}
	}

	return nil
}
