package repository

import (
	"context"
	"fmt"
	"time"

	domain "swallow-supplier/mongo/domain/yanolja"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *mongoRepository) UpsertItemIdDetails(ctx context.Context, itemIdDetails domain.ItemIdDetails) error {
	level.Info(r.logger).Log(
		"operation", "UpsertItemIdDetails",
		"itemIdDetails", itemIdDetails,
	)

	// Define the current timestamp
	currentTime := time.Now().UTC().Format(time.RFC3339)

	// Set `UpdatedAt` for all operations
	itemIdDetails.UpdatedAt = currentTime

	// If inserting, set `CreatedAt`
	if itemIdDetails.Id == "" {
		itemIdDetails.Id = primitive.NewObjectID().Hex()
		itemIdDetails.CreatedAt = currentTime
	}

	// Define the filter to find the document by OrderId
	filter := bson.M{"orderId": itemIdDetails.OrderId}

	// Define the update document
	update := bson.M{
		"$set": bson.M{
			"orderId":   itemIdDetails.OrderId,
			"items":     itemIdDetails.Items,
			"updatedAt": itemIdDetails.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"createdAt": itemIdDetails.CreatedAt,
		},
	}

	// Specify the upsert option (insert if not found)
	opts := options.Update().SetUpsert(true)

	// Perform the update operation
	collection := r.db.Collection("itemIdDetails")
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return fmt.Errorf("failed to upsert ItemIdDetails: %w", err)
	}
	return nil
}

// fetching all ItemIdDetail for syncing mongo to redis
func (r *mongoRepository) GetAllItemIdDetail(ctx context.Context) (itemIdDetail []domain.ItemIdDetails, err error) {
	filter := bson.M{} // Empty filter to match all documents

	collection := r.db.Collection("itemIdDetails")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var itemsiddetail domain.ItemIdDetails
		if err := cursor.Decode(&itemsiddetail); err != nil {
			return nil, err
		}
		itemIdDetail = append(itemIdDetail, itemsiddetail)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return itemIdDetail, nil
}
