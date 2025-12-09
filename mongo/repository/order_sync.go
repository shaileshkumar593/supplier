package repository

import (
	"context"

	domain "swallow-supplier/mongo/domain/yanolja"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRecentOrders retrieves orders updated/inserted in the last 5 minutes in UTC.
func (r *mongoRepository) GetRecentOrders(ctx context.Context) ([]domain.Model, error) {
	level.Info(r.logger).Log("repository method ", "GetRecentOrders")

	// Calculate 10 minutes ago in UTC.
	utcNow := time.Now().UTC()
	tenMinutesAgo := utcNow.Add(-11100 * time.Minute)

	// Define the filter.
	filter := bson.M{
		"updatedAt": bson.M{"$gte": tenMinutesAgo.Format(time.RFC3339)},
	}

	// Define the options (e.g., sorting by updatedAt descending).
	options := options.Find()
	options.SetSort(bson.D{{"updatedAt", -1}})

	collection := r.db.Collection("orders")
	// Execute the query.
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results.
	var orders []domain.Model
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
