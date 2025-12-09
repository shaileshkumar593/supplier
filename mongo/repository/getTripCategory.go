package repository

import (
	"context"
	"fmt"

	"swallow-supplier/request_response/trip"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetTripsByCategory retrieves trips based on multiple category combinations
func (r *mongoRepository) GetTripsByCategory(ctx context.Context, filters []trip.CategoryFilter) ([]string, error) {
	level.Info(r.logger).Log("repository method", "GetTripsByCategory")

	collection := r.db.Collection("category_mapping")

	// Construct OR filter query
	var orFilters []bson.M
	for _, f := range filters {
		orFilters = append(orFilters, bson.M{
			"yanolja_category_code":  f.CategoryCode,
			"yanolja_category_level": f.CategoryLevel,
		})
	}

	// Default to empty slice to avoid returning nil
	trips := []string{}

	// If no filters provided, return empty slice immediately
	if len(orFilters) == 0 {
		return trips, nil
	}

	// MongoDB filter
	filter := bson.M{"$or": orFilters}

	// Projection: Fetch only "trip_category_name"
	projection := options.Find().SetProjection(bson.M{"trip_category_name": 1, "_id": 0})

	// Query database
	cursor, err := collection.Find(ctx, filter, projection)
	if err != nil {
		return trips, fmt.Errorf("failed to fetch trips: %w", err)
	}
	defer cursor.Close(ctx)

	// Parse results
	var results []struct {
		TripCategoryName string `bson:"trip_category_name"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return trips, fmt.Errorf("failed to decode trips: %w", err)
	}

	// Extract "trip_category_name" values into the trips slice
	for _, result := range results {
		trips = append(trips, result.TripCategoryName)
	}

	return trips, nil
}
