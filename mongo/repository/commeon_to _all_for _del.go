package repository

import (
	"context"
	"fmt"
	"log"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
)

// DeleteAllDataAndCache removes all documents from collections and clears all Redis cache.
func (r *mongoRepository) DeleteAllIfNotEmpty(ctx context.Context) (map[string]int64, error) {
	//collections := []string{"plus", "product_images", "productview", "products"}
	collections := []string{"product_images", "productview", "products"}
	deletedCounts := make(map[string]int64)

	// Step 1: Delete data from MongoDB collections
	for _, collectionName := range collections {
		collection := r.db.Collection(collectionName)

		// Count documents before deletion
		count, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Printf("Error counting documents in %s: %v", collectionName, err)
			return nil, fmt.Errorf("failed to count documents in %s: %w", collectionName, err)
		}

		if count > 0 {
			// Delete all documents
			result, err := collection.DeleteMany(ctx, bson.M{})
			if err != nil {
				log.Printf("Error deleting records from %s: %v", collectionName, err)
				return nil, fmt.Errorf("failed to delete records from %s: %w", collectionName, err)
			}

			deletedCounts[collectionName] = result.DeletedCount
			log.Printf("Deleted %d records from %s", result.DeletedCount, collectionName)
		} else {
			deletedCounts[collectionName] = 0
			log.Printf("No records found in %s, skipping deletion", collectionName)
		}
	}

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(r.logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))

		return deletedCounts, err
	}
	// Step 2: Delete all cache entries from Redis
	keys, err := cacheLayer.Keys(ctx, "*")
	if err != nil {
		log.Printf("Error fetching Redis keys: %v", err)
		return nil, fmt.Errorf("failed to fetch Redis keys: %w", err)
	}

	cacheDeletedCount := int64(len(keys))
	if cacheDeletedCount > 0 {
		err = cacheLayer.Delete(ctx, keys)
		if err != nil {
			log.Printf("Error clearing Redis cache: %v", err)
			return nil, fmt.Errorf("failed to clear Redis cache: %w", err)
		}
		log.Printf("Deleted %d cache entries from Redis", cacheDeletedCount)
	} else {
		log.Println("No cache entries found, skipping deletion")
	}

	// Store cache deletion count in the result map
	deletedCounts["cache"] = cacheDeletedCount

	log.Println("Data deletion and full cache clearing process completed successfully")
	return deletedCounts, nil
}
