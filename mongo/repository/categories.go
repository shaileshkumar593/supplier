package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"swallow-supplier/mongo/domain/yanolja"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertCategories insert many categories at a time
func (r *mongoRepository) InsertCategories(ctx context.Context, categories []yanolja.Category) (err error) {
	level.Info(r.logger).Log("repo-method", "InsertCategories")

	// Correct the loop to avoid out-of-range error and assign a new ObjectID to each category.
	for i := 0; i < len(categories); i++ {
		categories[i].Id = primitive.NewObjectID().Hex()
	}

	// Get the collection reference.
	collection := r.db.Collection("categories")

	// Prepare the documents to be inserted.
	docs := make([]interface{}, len(categories))
	for i, category := range categories {
		docs[i] = category
	}

	// Perform the bulk insertion.
	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		level.Error(r.logger).Log("category insertions error", err)
		return err
	}

	return nil
}

// UpdateCategoriesByCategoryId update the category document based on categoryId
func (r *mongoRepository) UpdateCategoriesByCategoryId(ctx context.Context, categoryId int64, update map[string]any) (Id string, err error) {
	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}
	if len(update) > 0 {

		for key, val := range update {
			updateBson["$set"].(bson.M)[key] = val
		}
	}

	collection := r.db.Collection("categories")

	filter := bson.M{"categoryId": categoryId}

	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found
	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil {
		return "", fmt.Errorf("failed to update document: %w", err)
	}
	return result.UpsertedID.(string), nil
}

// fetch all the category record
func (r *mongoRepository) FindCategory(ctx context.Context) (record []yanolja.Category, err error) {
	filter := bson.M{} // Empty filter to match all documents

	collection := r.db.Collection("categories")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var category yanolja.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		record = append(record, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return record, nil
}

// FindCategoryByCategoryId filter category based on categoryId
func (r *mongoRepository) FindCategoryByCategoryId(ctx context.Context, categoryId int64) (record yanolja.Category, err error) {

	filter := bson.M{"categoryId": categoryId}

	collection := r.db.Collection("categories")

	err = collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return record, fmt.Errorf("no record found with categoryId %d", categoryId)
		}
		return record, err
	}
	return record, nil
}

// DeleteOrderByOrderId  delete the record by orderid
func (r *mongoRepository) DeleteCategoryByCategoryId(ctx context.Context, categoryId int64) (id string, err error) {
	collection := r.db.Collection("categories")

	filter := bson.M{"categoryId": categoryId}

	// Use FindOneAndDelete to find and delete the document
	var deletedDoc yanolja.Category
	err = collection.FindOneAndDelete(ctx, filter).Decode(&deletedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no document found with the given ID %d", categoryId)
		}
		return "", fmt.Errorf("failed to delete document with id %d with error %v", categoryId, err)
	}

	return deletedDoc.Id, nil
}

// BulkUpsertCategoryMapping  to update the category mapping to mongo
func (r *mongoRepository) BulkUpsertCategoryMapping(ctx context.Context, records []map[string]any) (map[string]int64, error) {
	level.Info(r.logger).Log("repository method", "BulkUpsertCategoryMapping", "record_count", len(records))

	countmap := make(map[string]int64)

	if len(records) == 0 {
		return nil, errors.New("no records to upsert")
	}

	collection := r.db.Collection("category_mapping")
	now := time.Now().UTC().Format(time.RFC3339)

	var operations []mongo.WriteModel

	for _, record := range records {
		// Normalize keys to snake_case before processing
		normalizedRecord := make(map[string]any)
		for key, val := range record {
			normalizedRecord[toSnakeCase(key)] = val
		}

		// Ensure required fields exist and are integers
		var categoryCode, categoryLevel int
		ok1, ok2 := false, false

		if v, exists := normalizedRecord["yanolja_category_code"]; exists {
			categoryCode, ok1 = convertToInt(v)
		}
		if v, exists := normalizedRecord["yanolja_category_level"]; exists {
			categoryLevel, ok2 = convertToInt(v)
		}

		if !ok1 || !ok2 {
			level.Warn(r.logger).Log("msg", "Skipping record due to missing required fields", "record", record)
			continue
		}

		// Query filter for upsert
		filter := bson.M{
			"yanolja_category_code":  categoryCode,
			"yanolja_category_level": categoryLevel,
		}

		// Prepare update document
		updateFields := bson.M{"updated_at": now}
		setOnInsertFields := bson.M{"created_at": now}

		// Convert record fields while keeping numeric values as int
		for key, val := range normalizedRecord {
			if key == "yanolja_category_code" || key == "yanolja_category_level" {
				continue // Already processed
			}
			switch v := val.(type) {
			case float64:
				// Convert float64 to int if it represents a whole number
				if v == float64(int(v)) {
					updateFields[key] = int(v)
				} else {
					updateFields[key] = v
				}
			case int, int32, int64:
				updateFields[key] = v
			case string:
				updateFields[key] = v
			default:
				updateFields[key] = v
			}
		}

		// Add bulk operation
		updateBson := bson.M{
			"$set":         updateFields,
			"$setOnInsert": setOnInsertFields,
		}

		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateBson).SetUpsert(true)
		operations = append(operations, updateModel)
	}

	// Execute bulk write
	if len(operations) > 0 {
		bulkOptions := options.BulkWrite().SetOrdered(false)
		result, err := collection.BulkWrite(ctx, operations, bulkOptions)
		if err != nil {
			level.Error(r.logger).Log("msg", "Bulk upsert failed", "error", err)
			return nil, fmt.Errorf("bulk upsert failed: %w", err)
		}

		level.Info(r.logger).Log(
			"msg", "Bulk upsert completed",
			"matched", result.MatchedCount,
			"modified", result.ModifiedCount,
			"inserted", result.UpsertedCount,
		)

		countmap["modified"] = result.ModifiedCount
		countmap["upserted"] = result.UpsertedCount

	} else {
		level.Warn(r.logger).Log("msg", "No valid records to upsert")
	}

	level.Info(r.logger).Log("msg", "Successfully inserted/upserted records")
	return countmap, nil
}

// function to convert to int
func convertToInt(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float64:
		if v == float64(int(v)) {
			return int(v), true
		}
		return 0, false // Reject non-whole numbers
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			return num, true
		}
	}
	return 0, false
}

// toSnakeCase converts a string from camelCase or space-separated format to snake_case.
func toSnakeCase(name string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	snake := re.ReplaceAllString(name, "${1}_${2}")
	snake = strings.ReplaceAll(snake, " ", "_")
	return strings.ToLower(snake)
}
