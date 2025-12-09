package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	svc "swallow-supplier/iface"
	"time"

	"github.com/go-kit/log"

	"github.com/go-kit/log/level"
)

// SyncItemIdDetail sync itemId detail to redis
func SyncItemIdDetail(ctx context.Context, logger log.Logger, mrepo svc.MongoRepository) error {
	/* level.Info(logger).Log(
		"method name ", "SyncItemIdDetail",
	) */

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		return err
	}

	itemIdAry, err := mrepo.GetAllItemIdDetail(ctx)
	if err != nil {
		return err
	}
	// Check if the Redis database is empty

	for _, item := range itemIdAry {
		// Convert the slice to JSON string
		itemAryJSON, err := json.Marshal(item)
		if err != nil {
			level.Error(logger).Log("msg", "Error marshaling  Item array", "error", err)
			return fmt.Errorf("Error marshaling  Item array: %w", err)
		}

		// Convert int64 to string
		orderIdstr := strconv.FormatInt(item.OrderId, 10)

		// Insert or update the record in Redis
		if _, err := cacheLayer.SetNX(ctx, orderIdstr, string(itemAryJSON), 7776000*time.Second); err != nil {
			level.Error(logger).Log(
				"msg", "Failed to set key in Redis",
				"key", orderIdstr,
				"value", string(itemAryJSON), // Log the serialized JSON for debugging
				"error", err,
			)
			return fmt.Errorf("failed to set key %s in Redis: %w", orderIdstr, err)
		}

		//level.Info(logger).Log("msg", "Successfully updated/inserted key in Redis", "key", orderIdstr)

	}

	return nil
}
