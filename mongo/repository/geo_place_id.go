package repository

import (
	"context"
	"fmt"
	"swallow-supplier/mongo/domain/trip"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *mongoRepository) BulkUpsertGooglePlaceIdOfProduct(ctx context.Context, docs []trip.GooglePlaceIdOfProduct) error {
	level.Info(r.logger).Log("operation", "BulkUpsertGooglePlaceIdOfProduct", "documentCount", len(docs))

	// Define current timestamp in RFC3339 format.
	currentTime := time.Now().UTC().Format(time.RFC3339)

	var models []mongo.WriteModel
	for _, doc := range docs {
		// Always set UpdatedAt.
		doc.UpdatedAt = currentTime

		// If inserting a new document, set ID and CreatedAt.
		if doc.Id == "" {
			doc.Id = primitive.NewObjectID().Hex()
			doc.CreatedAt = currentTime
		}

		// Filter by the unique ProductID.
		filter := bson.M{"productId": doc.ProductID}

		// Build the update document.
		update := bson.M{
			"$set": bson.M{
				"latitude":  doc.Latitude,
				"longitude": doc.Longitude,
				"placeId":   doc.PlaceId,
				"updatedAt": doc.UpdatedAt,
				"productId": doc.ProductID, // ensure ProductID is always stored
			},
			"$setOnInsert": bson.M{
				"createdAt": doc.CreatedAt,
			},
		}

		// Create an update model with upsert enabled.
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		models = append(models, model)
	}

	// Get a reference to the "geo_place_id" collection.
	collection := r.db.Collection("geo_place_id")
	opts := options.BulkWrite().SetOrdered(false)

	// Execute the bulk write.
	result, err := collection.BulkWrite(ctx, models, opts)
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return fmt.Errorf("failed to bulk upsert GooglePlaceIdOfProduct: %w", err)
	}

	level.Info(r.logger).Log("bulkUpsertResult", fmt.Sprintf("%+v", result))
	return nil
}
