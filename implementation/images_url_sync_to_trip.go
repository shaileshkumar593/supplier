package implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"swallow-supplier/config"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"
	tripservice "swallow-supplier/services/distributors/trip"

	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetImageSyncToTrip for syncing image url to trip
func (s *service) GetImageSyncToTrip(ctx context.Context) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	mrepo := s.mongoRepository[config.Instance().MongoDBName]
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

	syncingUrl, err := s.mongoRepository[config.Instance().MongoDBName].GetUnsyncedImagesForTrip(ctx)
	if err != nil {
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}
	imageUrlSync := trip.ImageSyncToTripRequest{
		Message: "ImageUrl",
		Data:    syncingUrl,
	}

	level.Info(logger).Log("imageUrl ", fmt.Sprintln(imageUrlSync))

	var tripsvc, _ = tripservice.New(ctx)
	resp, err = tripsvc.ImageUrlSyncToTrip(ctx, imageUrlSync)
	if err != nil {
		level.Error(logger).Log("error ", "Trip image sync ")
		return resp, err
	}

	imageUrlResp := resp.Body.([]map[string]interface{})
	syncingUrl = syncingUrl[:0]

	for _, respval := range imageUrlResp {

		update := map[string]any{
			"_id":       respval["id"],
			"productId": respval["productId"],
			"url":       respval["url"],
			"status":    respval["status"],
			"imageId":   respval["imageId"],
		}

		// Convert map to JSON bytes
		jsonData, err := json.Marshal(update)
		if err != nil {
			level.Error(logger).Log("error ", fmt.Sprintln("Error marshaling JSON:", err))
			return resp, err
		}

		// Convert JSON to struct
		var imageSync trip.ImageSyncToTrip
		err = json.Unmarshal(jsonData, &imageSync)
		if err != nil {
			level.Error(logger).Log("error ", fmt.Sprintln("Error unmarshaling JSON:", err))
			return resp, err
		}

		syncingUrl = append(syncingUrl, imageSync)

	}

	err = mrepo.BulkUpdateImageSyncStatus(ctx, syncingUrl)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintln("repository error for BulkUpdateImageSyncStatus "))
		return resp, err
	}

	// write code to add imagewithId  to mongo database as seprate collection
	resp.Code = "200"
	resp.Body = fmt.Sprintln("status updated successfully")
	return resp, nil
}

// UpdateTripImageSyncStatus for updating synced url status to trip
func (s *service) UpdateTripImageSyncStatus(ctx context.Context, req []trip.ImageSyncToTrip) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "UpdateTripImageSyncStatus",
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

	err = s.mongoRepository[config.Instance().MongoDBName].BulkUpdateImageSyncStatus(ctx, req)
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
