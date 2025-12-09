package repository

import (
	"context"
	"fmt"
	"swallow-supplier/request_response/trip"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Trip image sync
// GetUnsyncedImagesForTrip  get all unsync image url sync to trip for content api
func (r *mongoRepository) GetUnsyncedImagesForTrip(ctx context.Context) (images []trip.ImageSyncToTrip, err error) {
	level.Info(r.logger).Log("operation", "GetUnsyncedImagesForTrip")

	collection := r.db.Collection("product_images")

	// Define filter: Get records where `status` is NOT "Sync"
	filter := bson.M{
		"status": bson.M{"$ne": "Sync"},
	}

	// Define projection: Select only required fields
	projection := bson.M{
		"_id":       1,
		"productId": 1,
		"url":       1,
		"status":    1,
	}

	// Execute the query with projection
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return nil, fmt.Errorf("failed to fetch unsynced product images: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	if err := cursor.All(ctx, &images); err != nil {
		level.Error(r.logger).Log("decode error", err.Error())
		return nil, fmt.Errorf("failed to decode unsynced product images: %w", err)
	}

	level.Info(r.logger).Log("unsynced_images_count", len(images))
	return images, nil
}

// BulkUpdateImageSyncStatus update the status of sync image url
func (r *mongoRepository) BulkUpdateImageSyncStatus(ctx context.Context, images []trip.ImageSyncToTrip) error {
	level.Info(r.logger).Log("operation", "BulkUpdateImageSyncStatus", "image_count", len(images))

	if len(images) == 0 {
		level.Warn(r.logger).Log("message", "No images provided for update")
		return nil // No updates to process
	}

	collection := r.db.Collection("product_images")
	currentTime := time.Now().UTC().Format(time.RFC3339)

	var operations []mongo.WriteModel
	batchSize := 500 // Prevents exceeding MongoDB bulk write limits

	for _, image := range images {
		// Filter to match specific product image
		filter := bson.M{
			"productId": image.ProductId,
			"id":        image.ID,
			"url":       image.Url,
		}

		// Update `status` to "Sync" and set `updatedAt`
		update := bson.M{
			"$set": bson.M{
				"status":    "Sync",
				"imageId":   image.ImageId,
				"updatedAt": currentTime,
			},
		}

		operations = append(operations, mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update))

		// Execute batch updates
		if len(operations) >= batchSize {
			_, err := collection.BulkWrite(ctx, operations, options.BulkWrite().SetOrdered(false))
			if err != nil {
				level.Error(r.logger).Log("database error", err.Error())
				return fmt.Errorf("failed batch update: %w", err)
			}
			operations = nil // Reset batch
		}
	}

	// Execute remaining updates
	if len(operations) > 0 {
		_, err := collection.BulkWrite(ctx, operations, options.BulkWrite().SetOrdered(false))
		if err != nil {
			level.Error(r.logger).Log("database error", err.Error())
			return fmt.Errorf("failed final batch update: %w", err)
		}
	}

	level.Info(r.logger).Log("message", "Product image statuses updated to Sync")
	return nil
}

// for content api
// FetchTripImageIdsByProductID retrieves only the TripImageId array for a given ProductId.
func (r *mongoRepository) FetchTripImageIdsByProductID(ctx context.Context, productId int64) ([]string, error) {
	collection := r.db.Collection("product_images")

	filter := bson.M{"productId": productId}
	projection := bson.M{"imageId": 1, "_id": 0} // Include only tripImageId field

	var result struct {
		ImageIds []string `bson:"imageIds"`
	}

	// Using options.FindOne() to avoid undeclared error
	err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return empty if no document found
		}
		return nil, err
	}

	return result.ImageIds, nil
}
