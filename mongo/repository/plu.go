package repository

import (
	"context"
	"errors"
	"fmt"
	"swallow-supplier/mongo/domain/yanolja"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertAllPlu
func (r *mongoRepository) UpsertAllPlu(ctx context.Context, pluHash map[string]string) error {
	collection := r.db.Collection("plu_detail")

	var documentId string = "6510b8b8c9e77a3f4d6aab09"
	filter := bson.M{}

	if objectID, err := primitive.ObjectIDFromHex(documentId); err == nil {
		filter = bson.M{"_id": objectID} // If valid ObjectID, use it
	}

	// Update document with new records or insert if it doesn’t exist
	update := bson.M{
		"$set": bson.M{"pluhash": pluHash},
	}

	opts := options.Update().SetUpsert(true)

	// Perform upsert operation
	_, err1 := collection.UpdateMany(ctx, filter, update, opts)
	if err1 != nil {
		return fmt.Errorf("failed to upsert document: %w", err1)
	}

	return nil
}

// FindPluHashValue fetches the value corresponding to a given key in 'pluhash'
func (r *mongoRepository) FindPluHashValue(ctx context.Context, key string) (string, error) {
	collection := r.db.Collection("plus")

	// Define the document ID
	documentID := "6510b8b8c9e77a3f4d6aab09"
	objectID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {

		return "", fmt.Errorf("invalid document ID: %w", err)
	}

	// MongoDB aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match the document with the given _id
		{{Key: "$match", Value: bson.M{"_id": objectID}}},
		// Unwind the pluhash array
		{{Key: "$unwind", Value: "$pluhash"}},
		// Match the specific key inside the array
		{{Key: "$match", Value: bson.M{"pluhash." + key: bson.M{"$exists": true}}}},
		// Project only the required value
		{{Key: "$project", Value: bson.M{"_id": 0, "value": "$pluhash." + key}}},
	}

	// Run the aggregation query
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return "", fmt.Errorf("aggregation failed: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the result
	var result struct {
		Value string `bson:"value"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return "", fmt.Errorf("failed to decode result: %w", err)
		}
		return result.Value, nil
	}

	// If no document is found, return empty
	return "", nil
}

// FetchAllPluHashes
func (r *mongoRepository) FetchAllPluHashes(ctx context.Context) (map[string]string, error) {
	// level.Info(r.logger).Log("repository method", "FetchAllPluHashes")

	collection := r.db.Collection("productview")
	cursor, err := collection.Find(ctx, bson.M{}) // Fetch all documents
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	var allPluHashes = make(map[string]string)

	for cursor.Next(ctx) {
		var productview yanolja.ProductView
		if err := cursor.Decode(&productview); err != nil {
			return nil, fmt.Errorf("failed to decode productview: %w", err)
		}

		// Preallocate based on PluDetails length to optimize memory allocation
		for _, pluDetail := range productview.PluDetails {
			if len(pluDetail.PluHash) > 0 {
				for key, val := range pluDetail.PluHash {
					allPluHashes[key] = val
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return allPluHashes, nil
}

// filter plue value based on hash key
func (r *mongoRepository) FetchPluByKey(ctx context.Context, key string) (string, error) {
	collection := r.db.Collection("plus")

	// Get document ID dynamically
	documentID, err := r.GetDocumentID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get document ID: %w", err)
	}

	fmt.Println("===============  documentID ===================== ", documentID)
	// Define the filter
	filter := bson.M{"_id": documentID}

	// Define the projection to retrieve only the required key from "pluhash"
	projection := bson.M{"pluhash": 1}

	// Query MongoDB
	var result bson.M
	err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(r.logger).Log("error ", "no document found")
			return "", fmt.Errorf("no document found with _id: %s", documentID)
		}
		return "", err
	}

	// Extract the "pluhash" field
	pluhash, ok := result["pluhash"].(bson.M)
	if !ok {
		level.Error(r.logger).Log("error ", "pluhash field not found or invalid")
		return "", fmt.Errorf("pluhash field not found or invalid")
	}

	// Extract the value associated with the given key
	value, exists := pluhash[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in pluhash", key)
	}

	// Ensure the value is a string
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("invalid value type for key %s", key)
	}

	return strValue, nil
}

// Retrive id for accessing document
func (r *mongoRepository) GetDocumentID(ctx context.Context) (primitive.ObjectID, error) {
	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	collection := r.db.Collection("plus")

	err := collection.FindOne(ctx, bson.M{}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return primitive.NilObjectID, nil // No document found
		}
		return primitive.NilObjectID, err // Return other errors
	}

	return result.ID, nil
}

// FetchPluHashesByProductID
func (r *mongoRepository) FetchPluHashesByProductID(ctx context.Context, productId int64) (map[string]string, error) {
	// level.Info(r.logger).Log("repository method", "FetchPluHashesByProductID")

	collection := r.db.Collection("productview")

	// Query to filter by productId
	filter := bson.M{"productId": productId}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	allPluHashes := make(map[string]string)

	for cursor.Next(ctx) {
		var productview yanolja.ProductView
		if err := cursor.Decode(&productview); err != nil {
			return nil, fmt.Errorf("failed to decode productview: %w", err)
		}

		// Extract PLU hashes from PluDetails
		for _, pluDetail := range productview.PluDetails {
			for key, val := range pluDetail.PluHash {
				allPluHashes[key] = val
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	// If no records found, return an empty map instead of nil
	return allPluHashes, nil
}
