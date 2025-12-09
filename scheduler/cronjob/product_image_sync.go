package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	svc "swallow-supplier/iface"
	"swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/trip"
	tripservice "swallow-supplier/services/distributors/trip"
	"swallow-supplier/utils/constant"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
)

// ProductImagesForProcessing  processing image url to store on GCP
func ProductImagesForProcessing(ctx context.Context, logger log.Logger, mrepo svc.MongoRepository, products []yanolja.Product) (err error) {
	// Slice to store image URLs for processing
	var imageUrlForProcessing []yanolja.ImageUrlForProcessing

	// Map to keep track of unique combinations of ProductId, Url, and ImageTypeCode
	processed := make(map[string]bool)

	// Loop through each product image
	for _, product := range products {
		uniqueKey := fmt.Sprintf("%d", product.ProductID)

		// Skip if the combination has already been processed
		if processed[uniqueKey] {
			continue
		}
		for _, detail := range product.Images {
			for _, url := range detail.ImageURLs { // Loop through image URLs
				// Create the record and append it to the slice
				record := yanolja.ImageUrlForProcessing{
					ProductId:     product.ProductID,
					SupplierName:  product.SupplierName,
					Url:           url,
					ImageTypeCode: detail.ImageTypeCode,
					Status:        constant.STATUSNOTSYNC,
				}
				imageUrlForProcessing = append(imageUrlForProcessing, record)
			}

		}
		processed[uniqueKey] = true
		db := mrepo.GetMongoDb(ctx)
		productCollection := db.Collection("products")

		_, err = productCollection.UpdateOne(
			ctx,
			bson.M{"productId": product.ProductID, "imageScheduleStatus": false},
			bson.M{"$set": bson.M{"imageScheduleStatus": true}},
		)
		if err != nil {
			level.Error(logger).Log("error ", fmt.Sprintf("failed to update viewScheduleStatus for productId %d", product.ProductID))
			return fmt.Errorf("failed to update viewScheduleStatus for product: %w", err)
		}

	}

	// If no images were found, log and return early
	if len(imageUrlForProcessing) == 0 {
		level.Info(logger).Log("info ", "no image url available to insert")
		return nil
	}

	// Perform bulk insert for images to be processed
	err = mrepo.BulkInsertProductImagesUrl(ctx, imageUrlForProcessing)
	if err != nil {
		level.Error(logger).Log("error", "Error in BulkInsertProductImagesForProcessing", "error", err)
		return fmt.Errorf("Error in doing bulk upsert of images : %w", err)
	}

	return nil
}

// Get imageId from trip for the ImageUrl synced
func SyncImageUrlForImageId(ctx context.Context, logger log.Logger, mrepo svc.MongoRepository) (err error) {
	level.Info(logger).Log(
		"method name", "SyncImageUrlForImageId",
	)
	syncingUrl, err := mrepo.GetUnsyncedImagesForTrip(ctx)
	if err != nil {
		level.Error(logger).Log("error", "error in accesing GetUnsyncedImagesForTrip")
		return fmt.Errorf("Error in fetching all unsynced url : %w", err)
	}
	imageUrlSync := trip.ImageSyncToTripRequest{
		Message: "ImageUrl",
		Data:    syncingUrl,
	}

	// Trip sync
	level.Info(logger).Log("imageUrl ", fmt.Sprintln(imageUrlSync))

	var tripsvc, _ = tripservice.New(ctx)
	resp, err := tripsvc.ImageUrlSyncToTrip(ctx, imageUrlSync)
	if err != nil {
		level.Error(logger).Log("error ", "Trip image sync ")
		return err
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
			return err
		}

		// Convert JSON to struct
		var imageSync trip.ImageSyncToTrip
		err = json.Unmarshal(jsonData, &imageSync)
		if err != nil {
			level.Error(logger).Log("error ", fmt.Sprintln("Error unmarshaling JSON:", err))
			return err
		}

		syncingUrl = append(syncingUrl, imageSync)

	}

	err = mrepo.BulkUpdateImageSyncStatus(ctx, syncingUrl)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintln("repository error for BulkUpdateImageSyncStatus "))
		return err
	}

	return nil

}
