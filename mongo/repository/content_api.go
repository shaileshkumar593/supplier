package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"swallow-supplier/mongo/domain/trip"
	req_resp "swallow-supplier/request_response/trip"

	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BulkUpsertPackageContent  Handles bulk upsert operation for PackageContent
func (r *mongoRepository) BulkUpsertPackageContent(ctx context.Context, packageContents []trip.PackageContent) error {
	var operations []mongo.WriteModel
	now := time.Now().UTC().Format(time.RFC3339) // Get current time in ISO 8601 format

	for _, content := range packageContents {
		// Convert struct to BSON
		bsonDoc, err := bson.MarshalWithRegistry(bson.DefaultRegistry, content)
		if err != nil {
			return fmt.Errorf("failed to marshal package content: %w", err)
		}

		// Unmarshal to bson.M to modify fields
		var updateDoc bson.M
		if err := bson.Unmarshal(bsonDoc, &updateDoc); err != nil {
			return fmt.Errorf("failed to unmarshal package content: %w", err)
		}

		// Ensure `updatedAt` is always set correctly
		if val, exists := updateDoc["updatedAt"]; !exists || val == "" {
			updateDoc["updatedAt"] = now
		}

		// Remove `createdAt` from $set to avoid update conflicts
		delete(updateDoc, "createdAt")

		// Use `supplierProductId` as the filter for upsert
		filter := bson.M{"supplierProductId": content.SupplierProductId}

		// Define update operation
		update := bson.M{
			"$set": bson.M{
				"updatedAt": now, // Always update `updatedAt`
			},
			"$setOnInsert": bson.M{
				"createdAt": now, // Set `createdAt` only on insert
			},
		}

		// Merge remaining fields into `$set`
		for key, value := range updateDoc {
			update["$set"].(bson.M)[key] = value
		}

		// Handle `contractId` as MongoDB `NumberLong`
		if contractID, ok := update["$set"].(bson.M)["contractId"].(int64); ok {
			update["$set"].(bson.M)["contractId"] = bson.M{"$numberLong": fmt.Sprintf("%d", contractID)}
		}

		// Create upsert model
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true) // Insert if not found, update otherwise

		operations = append(operations, model)
	}

	// Execute bulk write
	opts := options.BulkWrite().SetOrdered(false)
	collection := r.db.Collection("package_content")
	result, err := collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		return fmt.Errorf("bulk write failed: %w", err)
	}

	log.Printf("Matched: %d, Modified: %d, Upserted: %d", result.MatchedCount, result.ModifiedCount, result.UpsertedCount)
	return nil
}

// BulkUpsertProductContent - Handles bulk upsert operation for ProductContent
func (r *mongoRepository) BulkUpsertProductContent(ctx context.Context, productContents []trip.ProuctContent) error {
	var operations []mongo.WriteModel
	now := time.Now().UTC().Format(time.RFC3339) // Current timestamp in ISO format

	for _, content := range productContents {
		// Convert struct to BSON
		bsonDoc, err := bson.MarshalWithRegistry(bson.DefaultRegistry, content)
		if err != nil {
			return fmt.Errorf("failed to marshal product content: %w", err)
		}

		// Convert BSON to bson.M for manipulation
		var updateDoc bson.M
		if err := bson.Unmarshal(bsonDoc, &updateDoc); err != nil {
			return fmt.Errorf("failed to unmarshal product content: %w", err)
		}

		// Remove "createdAt" from $set to prevent conflict
		delete(updateDoc, "createdAt")

		// Ensure `updatedAt` is explicitly set (if missing or empty)
		if val, exists := updateDoc["updatedAt"]; !exists || val == "" {
			updateDoc["updatedAt"] = now
		}

		// Ensure SupplierProductId is used as the filter
		filter := bson.M{"supplierProductId": content.SupplierProductId}

		// Define update operation
		update := bson.M{
			"$set": bson.M{
				"updatedAt": now, // Always update `updatedAt`
			},
			"$setOnInsert": bson.M{
				"createdAt": now, // Only set `createdAt` if inserting a new record
			},
		}

		// Merge remaining fields into `$set`
		for key, value := range updateDoc {
			update["$set"].(bson.M)[key] = value
		}

		// Ensure `contractId` is stored as a valid `NumberLong`
		if contractID, ok := update["$set"].(bson.M)["contractId"].(int64); ok {
			update["$set"].(bson.M)["contractId"] = bson.M{"$numberLong": fmt.Sprintf("%d", contractID)}
		}

		// MongoDB upsert model
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true) // Insert if not found, update otherwise

		operations = append(operations, model)
	}

	// Execute bulk write
	opts := options.BulkWrite().SetOrdered(false)
	collection := r.db.Collection("product_content")
	result, err := collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		return fmt.Errorf("bulk write failed: %w", err)
	}

	log.Printf("Matched: %d, Modified: %d, Upserted: %d", result.MatchedCount, result.ModifiedCount, result.UpsertedCount)
	return nil
}

// GetProductContentNotSync retrieves product content which is notSync
func (r *mongoRepository) GetProductContentNotSync(ctx context.Context) ([]trip.ProuctContent, error) {
	level.Info(r.logger).Log("operation", "GetProductContentNotSync")

	collection := r.db.Collection("product_content")

	// Define the filter to find documents with `notSync: true`
	filter := bson.M{"syncStatus": "NotSync"}

	// Query MongoDB without sorting or limiting
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		level.Error(r.logger).Log("error", "failed to fetch unsynchronized product content", "reason", err.Error())
		return nil, fmt.Errorf("failed to fetch unsynchronized product content: %w", err)
	}
	defer cursor.Close(ctx)

	var productContents []trip.ProuctContent
	if err := cursor.All(ctx, &productContents); err != nil {
		level.Error(r.logger).Log("error", "failed to decode product content", "reason", err.Error())
		return nil, fmt.Errorf("failed to decode product content: %w", err)
	}

	level.Info(r.logger).Log("operation", "GetProductContentNotSync", "count", len(productContents), "status", "success")
	return productContents, nil
}

// GetPackageContentNotSync  retrives package content notSync
func (r *mongoRepository) GetPackageContentNotSync(ctx context.Context) ([]trip.PackageContent, error) {
	level.Info(r.logger).Log("operation", "GetPackageContentNotSync")

	collection := r.db.Collection("package_content")

	// Define the filter to find documents where SyncStatus is "NotSync"
	filter := bson.M{"SyncStatus": "NotSync"}

	// Query MongoDB without sorting or limiting
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		level.Error(r.logger).Log("error", "failed to fetch unsynchronized package content", "reason", err.Error())
		return nil, fmt.Errorf("failed to fetch unsynchronized package content: %w", err)
	}
	defer cursor.Close(ctx)

	var packageContents []trip.PackageContent
	if err := cursor.All(ctx, &packageContents); err != nil {
		level.Error(r.logger).Log("error", "failed to decode package content", "reason", err.Error())
		return nil, fmt.Errorf("failed to decode package content: %w", err)
	}

	level.Info(r.logger).Log("operation", "GetPackageContentNotSync", "count", len(packageContents), "status", "success")
	return packageContents, nil
}

// BulkUpdateSyncStatus  to update the SyncStatus of the ProductContent  or packageContent
func (r *mongoRepository) BulkUpdateSyncStatus(ctx context.Context, updates req_resp.TripMessageForSync) error {
	level.Info(r.logger).Log("operation", "BulkUpdateSyncStatus")

	var operations []mongo.WriteModel

	now := time.Now().UTC().Format(time.RFC3339)

	for _, contentStatus := range updates.Status {
		if contentStatus.SupplierProductId == "" {
			continue
		}

		filter := bson.M{"supplierProductId": contentStatus.SupplierProductId}
		update := bson.M{
			"$set": bson.M{
				"syncStatus": contentStatus.SyncStatus,
				"updatedAt":  now,
			},
		}

		operations = append(operations, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	if len(operations) == 0 {
		return nil // No valid updates
	}

	var collection *mongo.Collection

	if strings.EqualFold(updates.Message, "product") {
		collection = r.db.Collection("product_content")

	} else if strings.EqualFold(updates.Message, "product") {
		collection = r.db.Collection("package_content")

	}
	_, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		level.Error(r.logger).Log("error", "failed to bulk update SyncStatus", "reason", err.Error())
		return fmt.Errorf("failed to bulk update sync status: %w", err)
	}

	level.Info(r.logger).Log("operation", "BulkUpdateSyncStatus", "status", "success")
	return nil
}
