package repository

import (
	"context"
	"fmt"
	"swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/utils/constant"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BulkInserProductImagesUrl   creating product_images on syncing url
func (r *mongoRepository) BulkInsertProductImagesUrl(ctx context.Context, images []yanolja.ImageUrlForProcessing) error {
	level.Info(r.logger).Log("operation", "BulkInsertProductImagesUrl", "image_count", len(images))

	collection := r.db.Collection("product_images")
	currentTime := time.Now().UTC().Format(time.RFC3339)

	var documents []interface{}
	var filters []bson.M

	for _, image := range images {
		filters = append(filters, bson.M{
			"productId":     image.ProductId,
			"url":           image.Url,
			"imageTypeCode": image.ImageTypeCode,
		})
	}

	existingCursor, err := collection.Find(ctx, bson.M{"$or": filters})
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return fmt.Errorf("failed to query existing records: %w", err)
	}

	existingRecords := make(map[string]struct{})
	for existingCursor.Next(ctx) {
		var existingDoc bson.M
		if err := existingCursor.Decode(&existingDoc); err != nil {
			level.Error(r.logger).Log("database error", err.Error())
			return fmt.Errorf("failed to decode existing record: %w", err)
		}
		key := fmt.Sprintf("%v|%v|%v", existingDoc["productId"], existingDoc["url"], existingDoc["imageTypeCode"])
		existingRecords[key] = struct{}{}
	}

	for _, image := range images {
		key := fmt.Sprintf("%v|%v|%v", image.ProductId, image.Url, image.ImageTypeCode)
		if _, exists := existingRecords[key]; !exists {
			documents = append(documents, bson.M{
				"productId":     image.ProductId,
				"supplierName":  image.SupplierName,
				"url":           image.Url,
				"imageTypeCode": image.ImageTypeCode,
				"status":        constant.STATUSNOTSYNC,
				"createdAt":     currentTime,
				"updatedAt":     currentTime,
			})
		}
	}

	if len(documents) > 0 {
		_, err := collection.InsertMany(ctx, documents)
		if err != nil {
			level.Error(r.logger).Log("database error", err.Error())
			return fmt.Errorf("failed bulk insertion operation: %w", err)
		}
		level.Info(r.logger).Log("message", "Bulk insert completed successfully", "inserted_count", len(documents))
	} else {
		level.Info(r.logger).Log("message", "No new records to insert")
	}

	return nil
}

// BulkUpdateProductImageStatusAndImageId for prcessing result
func (r *mongoRepository) BulkUpdateProductImageStatusAndImageId(ctx context.Context, images []yanolja.ImageUrlForProcessing) error {
	level.Info(r.logger).Log("operation", "BulkUpdateProductImageStatus", "image_count", len(images))

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
			"productId":     image.ProductId,
			"url":           image.Url,
			"imageTypeCode": image.ImageTypeCode,
		}

		// Update `status` to "Sync" and set `updatedAt`
		update := bson.M{
			"$set": bson.M{
				"status":    image.Status,
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

// GetUnsyncedProductImagesForProcessing
func (r *mongoRepository) GetUnsyncedProductImagesForProcessing(ctx context.Context) ([]yanolja.ImageUrlForProcessing, error) {
	level.Info(r.logger).Log("operation", "GetUnsyncedProductImages")

	collection := r.db.Collection("product_images")

	// Define filter: Get records where `status` is NOT "Sync"
	filter := bson.M{
		"status": bson.M{"$ne": "Sync"}, // $ne (not equal)
	}

	// Execute the query
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return nil, fmt.Errorf("failed to fetch unsynced product images: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var images []yanolja.ImageUrlForProcessing
	if err := cursor.All(ctx, &images); err != nil {
		level.Error(r.logger).Log("decode error", err.Error())
		return nil, fmt.Errorf("failed to decode unsynced product images: %w", err)
	}

	level.Info(r.logger).Log("unsynced_images_count", len(images))
	return images, nil
}
