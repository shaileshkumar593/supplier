package cronjob

import (
	"context"
	"fmt"

	//"swallow-supplier/caches/cache"
	//"swallow-supplier/config"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	svc "swallow-supplier/iface"

	//"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// SchedulePluUpsertToRedis update plu to redis for order creation
func SchedulePluUpsertToRedis(ctx context.Context, mrepo svc.MongoRepository, logger log.Logger) (err error) {
	allPluHashes, err := mrepo.FetchAllPluHashes(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "Error in FetchAllPluHashes", "error", err)
		return fmt.Errorf("error in FetchAllPluHashes: %w", err)
	}

	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("msg", "Error initializing cache layer", "error", err)
		return fmt.Errorf("error initializing cache layer: %w", err)
	}

	if len(allPluHashes) > 0 {

		for uid, plu := range allPluHashes {

			exist := cacheLayer.Exist(ctx, uid)
			if exist == 0 {

				// Insert or update the record in Redis
				if err := cacheLayer.Set(ctx, uid, plu); err != nil {
					level.Error(logger).Log(
						"msg", "Failed to set key in Redis for every second",
						"key", uid,
						"value", plu, // Log the serialized JSON for debugging
						"error", err,
					)
					return fmt.Errorf("failed to set key %s in Redis: %w", uid, err)
				}
			}

		}
	}
	return nil
}
