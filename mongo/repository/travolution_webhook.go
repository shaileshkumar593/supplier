package repository

import (
	"context"
	"fmt"
	"swallow-supplier/mongo/domain/travolution"
	domain "swallow-supplier/mongo/domain/travolution"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertWebhook Insert/update webhook to mongo
func (r *mongoRepository) UpsertWebhook(ctx context.Context, payload domain.Webhook, field string) (string, error) {
	collection := r.db.Collection("travolution_webhooks")

	// Unique key = referenceNumber + eventType
	filter := bson.M{
		"data.referenceNumber": payload.Data.ReferenceNumber,
		"data.orderNumber":     payload.Data.OrderNumber,
		"eventType":            payload.EventType,
	}

	now := time.Now().UTC().Format(time.RFC3339)

	// Prepare dynamic update
	updateFields := bson.M{
		"eventType": payload.EventType,
		"createdAt": payload.CreatedAt,
		"updatedAt": now,
	}

	update := bson.M{"$set": updateFields}

	opts := options.Update().SetUpsert(true)

	res, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "", err
	}

	// Return inserted or updated ID
	if res.UpsertedID != nil {
		if oid, ok := res.UpsertedID.(primitive.ObjectID); ok {
			return oid.Hex(), nil
		}
		return fmt.Sprintf("%v", res.UpsertedID), nil
	}

	return payload.ID, nil
}

// InsertWebhook inserts a new webhook record into the travolution_webhook collection.
func (r *mongoRepository) UpsertTravolutionWebhook(ctx context.Context, payload travolution.Webhook) (string, error) {
	collection := r.db.Collection("travolution_webhooks")
	now := time.Now().UTC().Format(time.RFC3339)

	// Assign ID if missing
	if payload.ID == "" {
		payload.ID = primitive.NewObjectID().Hex()
	}

	// Ensure timestamps
	if payload.CreatedAt == "" {
		payload.CreatedAt = now
	}
	payload.UpdatedAt = now

	// Upsert filter: uniquely identify a webhook by eventType, orderNumber, referenceNumber, dateAt
	filter := bson.M{
		"eventType":            payload.EventType,
		"data.orderNumber":     payload.Data.OrderNumber,
		"data.referenceNumber": payload.Data.ReferenceNumber,
		"data.dateAt":          payload.Data.DateAt,
	}

	// Update document
	update := bson.M{
		"$set": bson.M{
			"data":      payload.Data, // overwrite the latest payload data
			"updatedAt": payload.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":       payload.ID,
			"createdAt": now,
		},
	}

	opts := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upsert webhook: %w", err)
	}

	// Return the inserted or existing document ID
	if result.UpsertedID != nil {
		if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
			return oid.Hex(), nil
		}
		return payload.ID, nil
	}

	// Fetch existing doc ID if updated
	var existing travolution.Webhook
	if err := collection.FindOne(ctx, filter).Decode(&existing); err != nil {
		return "", fmt.Errorf("failed to fetch updated webhook: %w", err)
	}

	return existing.ID, nil
}

// UpsertWebhookToOrder   append webhook to particular order
func (r *mongoRepository) UpsertWebhookToOrder(
	ctx context.Context,
	payload domain.Webhook,
	update map[string]any,
) (string, error) {
	ordersColl := r.db.Collection("travolution_orders")
	now := time.Now().UTC().Format(time.RFC3339)

	// Create event entry
	eventEntry := domain.WebHookData{
		EventType: payload.EventType,
		DateAt:    payload.Data.DateAt,
	}

	// Order filter
	orderFilter := bson.M{
		"orderNumber":     payload.Data.OrderNumber,
		"referenceNumber": payload.Data.ReferenceNumber,
	}

	// Ensure update map
	if update == nil {
		update = make(map[string]any)
	}
	update["updatedAt"] = now

	// Build update document
	updateBson := bson.M{
		"$set": bson.M{},
		"$push": bson.M{
			"eventHistory": eventEntry,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}

	// Apply fields to $set
	for key, val := range update {
		updateBson["$set"].(bson.M)[key] = val
	}

	opts := options.Update().SetUpsert(true)

	res, err := ordersColl.UpdateOne(ctx, orderFilter, updateBson, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upsert order with eventHistory: %w", err)
	}

	// Return inserted/updated ID
	if res.UpsertedID != nil {
		if oid, ok := res.UpsertedID.(primitive.ObjectID); ok {
			return oid.Hex(), nil
		}
		return fmt.Sprintf("%v", res.UpsertedID), nil
	}

	// Updated existing document → return composite key
	return fmt.Sprintf("%s-%s", payload.Data.OrderNumber, payload.Data.ReferenceNumber), nil
}
