package repository

import (
	"context"
	"fmt"
	"swallow-supplier/mongo/domain/yanolja"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertRegions insert many Regions detail at a time
func (r *mongoRepository) InsertRegions(ctx context.Context, regions []yanolja.Region) (err error) {
	level.Info(r.logger).Log("repo-method", "InsertRegions")
	for i := 0; i < len(regions); i++ {
		regions[i].Id = primitive.NewObjectID().Hex()
	}

	// Define the index model for the unique constraint on CategoryCode
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "regionId", Value: 1}}, // Index on categoryCode in ascending order
		Options: options.Index().SetUnique(true),     // Set the index to be unique
	}

	// Create the unique index on the CategoryCode field if it doesn't already exist
	collection := r.db.Collection("regions")
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		level.Error(r.logger).Log("error creating unique index on CategoryCode", err)
		return err
	}

	docs := make([]interface{}, len(regions))
	for i, category := range regions {
		docs[i] = category
	}

	//collection := r.db.Collection("regions")
	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		level.Error(r.logger).Log("regional insertions error", err)
		return err
	}

	return nil
}

// UpdateRegionsByRegionId update the document by regiondId
func (r *mongoRepository) UpdateRegionsByRegionId(ctx context.Context, regiondId int64, update map[string]any) (Id string, err error) {
	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}
	if len(update) > 0 {

		for key, val := range update {
			updateBson["$set"].(bson.M)[key] = val
		}
	}

	collection := r.db.Collection("regions")

	filter := bson.M{"regionId": regiondId}

	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found
	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil {
		return "", fmt.Errorf("failed to update document: %w", err)
	}
	return result.UpsertedID.(string), nil
}

// FindRegion fetch all the regions record
func (r *mongoRepository) FindRegion(ctx context.Context) (records []yanolja.Region, err error) {
	filter := bson.M{} // Empty filter to match all documents

	collection := r.db.Collection("regions")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var region yanolja.Region
		if err := cursor.Decode(&region); err != nil {
			return nil, err
		}
		records = append(records, region)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// FindRecordByRegionId fetch the record by regionId
func (r *mongoRepository) FindRecordByRegionId(ctx context.Context, regionId int64) (record yanolja.Region, err error) {

	filter := bson.M{"regionId": regionId}

	collection := r.db.Collection("regions")

	err = collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return record, fmt.Errorf("no record found with regionId %d", regionId)
		}
		return record, err
	}
	return record, nil
}

// DeleteRegionByRegionId delete the record by RegionId
func (r *mongoRepository) DeleteRegionByRegionId(ctx context.Context, regionId int64) (id string, err error) {
	collection := r.db.Collection("regions")

	filter := bson.M{"regionId": regionId}

	// Use FindOneAndDelete to find and delete the document
	var deletedDoc yanolja.Region
	err = collection.FindOneAndDelete(ctx, filter).Decode(&deletedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no document found with the given ID %d", regionId)
		}
		return "", fmt.Errorf("failed to delete document with id %d with error %v", regionId, err)
	}

	return deletedDoc.Id, nil
}
