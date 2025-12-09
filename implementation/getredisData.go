package implementation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	customError "swallow-supplier/error"

	//domain "swallow-supplier/mongo/domain/yanolja"

	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/scheduler/cronjob"
	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
)

var PLUMap = map[string]string{

	"Ur0wLjG561oQWHO": "10141711|39|11538777||",
	"pxIFUphrJ0GTQ1b": "10141711|39|11538775||",
	"etxEiIGYBcAt9oH": "10141711|39|11538776||",
	"BpKsZpkHERIFCCC": "10177758|43|11961176||",
	"isRfQLGNIULrHpe": "10177758|43|11961177||",
	"h4G8gI59sq5fBK7": "10177758|43|11961178||",
}

// to manual test redis issue
func (s *service) GetRedisData(ctx context.Context) (resp yanolja.Response, err error) {

	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetRedisData",
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

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		resp.Code = "500"
		return resp, err
	}
	var keyary []string
	keyary, err = cacheLayer.Keys(ctx, "*")
	if err != nil {
		resp.Code = "500"
		return resp, err
	}

	redisdata := make(map[string]string, len(keyary))
	for _, key := range keyary {
		fmt.Println("--------------------------key----------------------- ", key)
		strd, err := cacheLayer.Get(ctx, key)
		fmt.Println("**************************************** : ", strd)
		if err != nil {
			fmt.Println("===============Error ===================  ", err)
			level.Error(logger).Log("error", fmt.Sprintf("Error in accessing get function of cache: %s", err))
			continue
		}
		redisdata[key] = strd
	}
	resp.Code = "200"
	resp.Body = redisdata
	return resp, nil

}

func (s *service) GetPluToRedis(ctx context.Context) (resp yanolja.Response, err error) {

	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetPluToRedis",
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

	// for testing key based value ofplu
	/* key := "mwF5e3sL0P6Fp2m"

	val, err := s.mongoRepository[config.Instance().MongoDBName].FetchPluByKey(ctx, key)
	if err != nil {
		return resp, err
	} */
	allPluHashes, err := s.mongoRepository[config.Instance().MongoDBName].FetchAllPluHashes(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "Error in FetchAllPluHashes", "error", err)
		return resp, fmt.Errorf("Error in FetchAllPluHashes: %w", err)
	}

	fmt.Println("************************ length of PluHash ************** : ", len(allPluHashes))

	err = s.mongoRepository[config.Instance().MongoDBName].UpsertAllPlu(ctx, allPluHashes)
	if err != nil {
		level.Error(logger).Log("error ", "allPluUpsert errror")
		return resp, err
	}

	key := "1lkpDphZYClNnoW"
	val, err := s.mongoRepository[config.Instance().MongoDBName].FindPluHashValue(ctx, key)

	if err != nil {
		return resp, err
	}
	fmt.Println("***************  key val *************** ", val)

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		resp.Code = "500"
		return resp, err
	}
	var keyary []string
	keyary, err = cacheLayer.Keys(ctx, "*")
	if err != nil {
		resp.Code = "500"
		return resp, err
	}

	fmt.Println("********************** keyAry ******************** ", keyary)

	var plumap = make(map[string]string)

	for _, key := range keyary {
		val, _ := cacheLayer.Get(ctx, key)
		fmt.Println("******************* val **********************", val)
		plumap[key] = val
	}

	fmt.Println("************************ plumap ********************", plumap, len(plumap))
	hash, _ := s.mongoRepository[config.Instance().MongoDBName].FetchPluHashesByProductID(ctx, 10015449)

	fmt.Println("&&&&&&&&&&&&&&&&&&& hash &&&&&&&&&&&&&&&&&&&&& ", len(hash), hash)
	resp.Code = "200"
	resp.Body = hash
	return resp, nil

}

func (s *service) DeleteRecoRdIfNotEmpty(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "DeleteRecoRdIfNotEmpty",
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

	countmap, err := s.mongoRepository[config.Instance().MongoDBName].DeleteAllIfNotEmpty(ctx)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf(" repository error while deleting record : %s", err))
		resp.Code = "500"
		return resp, err
	}
	resp.Code = "200"
	resp.Body = countmap

	return
}

// manual testing
func (s *service) MonitorProductUpdateSvc(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "MonitorProductUpdateSvc",
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

	err = cronjob.MonitorProductUpdates(ctx, s.mongoRepository[config.Instance().MongoDBName], logger)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf(" issue in running monitor product updates cron job : %s", err))
		resp.Code = "500"
		return resp, err
	}
	resp.Code = "200"
	resp.Body = "Successfully completed monitor product update"
	return
}

// for manual testing of sync plu to redis

func (s *service) SyncAllPluToRedis(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "MonitorProductUpdateSvc",
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

	err = cronjob.SchedulePluUpsertToRedis(ctx, s.mongoRepository[config.Instance().MongoDBName], logger)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf(" issue in running monitor product updates cron job : %s", err))
		resp.Code = "500"
		return resp, err
	}
	resp.Code = "200"
	resp.Body = "Successfully completed update to redis"
	return
}

// for manual testing
func (s *service) DeleteRedisData(ctx context.Context) (resp yanolja.Response, err error) {

	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "DeleteRedisData",
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

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		resp.Code = "500"
		return resp, err
	}
	var keyary []string
	keyary, err = cacheLayer.Keys(ctx, "*")
	if err != nil {
		resp.Code = "500"
		return resp, err
	}

	err = cacheLayer.Delete(ctx, keyary)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error in deleting record in cache layer: %s", err))
		resp.Code = "500"
		return resp, err
	}
	resp.Code = "200"
	resp.Body = keyary
	return resp, nil

}

func (s *service) FindRedisKeyValue(ctx context.Context, key string) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetRedisData",
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

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		resp.Code = "500"
		return resp, err
	}

	strd, err := cacheLayer.Get(ctx, key)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error in accessing get function of cache: %s", err))
		resp.Body = fmt.Sprintln("error in accessing the get function")
		resp.Code = "500"
		return resp, err
	}

	resp.Body = strd
	resp.Code = "200"
	return resp, nil

}

// write plu to redis
func (s *service) UpdatePluToRedis(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "UpdatePluToRedis",
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

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		resp.Code = "500"
		return resp, err
	}

	for key, plu := range PLUMap {
		// Insert or update the record in Redis
		if err := cacheLayer.Set(ctx, key, plu); err != nil {
			level.Error(logger).Log(
				"msg", "Failed to set key in Redis",
				"key", key,
				"value", plu, // Log the serialized JSON for debugging
				"error", err,
			)
			resp.Body = fmt.Sprintln("msg", "Failed to set key in Redis")
			resp.Code = "500"
			return resp, fmt.Errorf("failed to set key %s in Redis: %w", key, err)
		}
	}

	resp.Body = fmt.Sprintln("successfully updated the redis ")
	resp.Code = "200"
	return resp, nil
}

// GetPluFromRedis get plu from redis based on productId, Version and VariantId
func (s *service) GetPluFromRedis(ctx context.Context, req yanolja.PluRequest) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetPluFromRedis",
		"Request ID", requestID,
	)

	var pluObj = make(map[string]string, 0)
	var hash string

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	productview, err := s.mongoRepository[config.Instance().MongoDBName].GetPLUDetails(ctx, req.ProductId)
	if err == mongo.ErrNoDocuments {
		level.Error(logger).Log("repository error", "no record exist based on productId  ", err)
		resp.Code = "404"
		resp.Body = fmt.Sprintf("no productvie exist with productId %d", req.ProductId)
		return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetPluFromRedis")
	}

	if productview.PluDetails != nil {
		for _, detail := range productview.PluDetails {
			for key, plu := range detail.PluHash {
				hash = key
				fmt.Println("  key ", key, " plu ", plu)
				level.Info(logger).Log("plu ", plu)
				detail := strings.Split(plu, "|") // ProductID-ProductVersion-VariantID

				variantId, err := strconv.ParseInt(detail[2], 10, 64)
				if err != nil {
					level.Error(logger).Log("error :", err)
					resp.Code = "500"
					resp.Body = "data conversion error"
					return resp, customError.NewErrorCustom(ctx, resp.Code, resp.Body.(string), "", http.StatusInternalServerError, "GetPluFromRedis")
				}
				fmt.Println(" variant ", variantId, " req.VarinatId ", req.VariantId)

				if req.VariantId == variantId {
					pluObj[hash] = plu
					break
				} else {
					hash = ""
					continue
				}

			}

		}
	}

	if hash == "" {
		level.Error(logger).Log("error", fmt.Sprintf("no plu exist for productid %d variantId %d and version %d", req.ProductId, req.VariantId, req.ProductVersion), err)
		resp.Code = "404"
		resp.Body = fmt.Sprintf("no plu exist for requested variant %d", req.VariantId)
		return resp, customError.NewErrorCustom(ctx, resp.Code, resp.Body.(string), "", http.StatusNotFound, "GetPluFromRedis")
	}
	fmt.Println(" ***************  plu map ", pluObj)
	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		resp.Code = "500"
		resp.Body = "Error initializing cache layer"
		return resp, customError.NewErrorCustom(ctx, resp.Code, resp.Body.(string), "", http.StatusNotFound, "GetPluFromRedis")

	}

	pluRecord, err := cacheLayer.Get(ctx, hash)
	if err != nil {
		err1 := cacheLayer.Set(ctx, hash, pluObj[hash]) // to gurantee that plu is avail in redis
		if err1 != nil {
			level.Error(logger).Log(
				"msg", "Failed to set key in Redis",
				"key", hash,
				"value", pluObj[hash], // Log the serialized JSON for debugging
				"error", err,
			)
			resp.Body = fmt.Sprintln("msg", "Failed to set key in Redis")
			resp.Code = "500"
			resp.Body = fmt.Errorf("failed to set key %s in Redis: %w", hash, err)
			return resp, customError.NewErrorCustom(ctx, resp.Code, resp.Body.(string), "", http.StatusNotFound, "GetPluFromRedis")
		}

	}

	fmt.Println("<<<<<<<<<<<<<<<<<  pluRecord >>>>>>>>>>>>>>>>>>>> ", pluRecord)
	pluObj[hash] = pluRecord
	resp.Body = pluObj
	resp.Code = "200"
	return resp, nil
}
