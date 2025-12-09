package repository

import (
	"context"
	"fmt"
	"swallow-supplier/mongo/domain/yanolja"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertProduct
func (r *mongoRepository) UpsertTripSequencdId(ctx context.Context, mapper yanolja.SequenceIdDetail) error {
	level.Info(r.logger).Log(
		"operation", "UpsertTripSequencdId",
		"partnerOrderId", mapper.PartnerOrderID,
	)

	// Define the current timestamp
	currentTime := time.Now().UTC().Format(time.RFC3339)

	// Set `UpdatedAt` for all operations
	mapper.UpdatedAt = currentTime

	// If inserting, set `CreatedAt`
	if mapper.Id == "" {
		mapper.Id = primitive.NewObjectID().Hex()
		mapper.CreatedAt = currentTime
	}

	// Define the filter to find the product by ProductID
	filter := bson.M{"partnerOrderId": mapper.PartnerOrderID}

	// Define the update document
	update := bson.M{
		"$set": bson.M{
			"sequenceId":          mapper.SequenceId,
			"sequenceIdGen":       mapper.SequenceIdGen,
			"orderId":             mapper.OrderId,
			"partnerOrderGroupID": mapper.PartnerOrderGroupID,
			"partnerOrderId":      mapper.PartnerOrderID,
			"updatedAt":           mapper.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"createdAt": mapper.CreatedAt,
		},
	}

	// Specify the upsert option (insert if not found)
	opts := options.Update().SetUpsert(true)

	// Perform the update operation
	collection := r.db.Collection("SequenceIdDetail")
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return fmt.Errorf("failed to upsert product: %w", err)
	}

	return nil
}
