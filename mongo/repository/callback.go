package repository

import (
	"context"
	"fmt"
	"time"

	"swallow-supplier/mongo/domain/trip"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertCallBackDetail   insert callback record to mongo
func (r *mongoRepository) InsertCallBackDetail(ctx context.Context, callbackDetail trip.CallBackDetail) (id string, err error) {
	level.Info(r.logger).Log("repo-method", "InsertCallBackDetail")

	collection := r.db.Collection("callBack_detail")

	// Set initial status and formatted timestamps
	currentTime := time.Now().UTC().Format(time.RFC3339)
	callbackDetail.GGTToChannelStatus = "PROCESSING"
	callbackDetail.ChannelToGGTStatus = ""
	createdAt := currentTime
	updatedAt := currentTime
	callbackDetail.Id = primitive.NewObjectID().Hex()

	// Convert struct to BSON with timestamps
	document := bson.M{
		"_id":                callbackDetail.Id,
		"tripCallBackInfo":   callbackDetail.ChannelCallBackInfo,
		"ggtToChannelStatus": callbackDetail.GGTToChannelStatus,
		"channelToGGTStatus": callbackDetail.ChannelToGGTStatus,
		"createdAt":          createdAt,
		"updatedAt":          updatedAt,
	}

	_, err = collection.InsertOne(ctx, document)
	if err != nil {
		level.Error(r.logger).Log("error", "Failed to insert callback detail", "id", callbackDetail.Id, "err", err)
		return "", err
	}

	level.Info(r.logger).Log("info", "Inserted callback detail successfully", "id", callbackDetail.Id)
	return callbackDetail.Id, nil
}

// UpdateCallBackStatus   update status
func (r *mongoRepository) UpdateCallBackStatus(ctx context.Context, id, status string) error {
	collection := r.db.Collection("callBack_detail")

	// Validate status input
	validStatuses := map[string]bool{
		"PROCESSING": true,
		"FAILED":     true,
		"SUCCESS":    true,
	}
	if !validStatuses[status] {
		level.Error(r.logger).Log("error", "Invalid status provided", "id", id, "status", status)
		return fmt.Errorf("invalid status: %s", status)
	}

	// Prepare update filter and update fields
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"ggtToChannelStatus": status,
			"channelToGGTStatus": status,
			"updatedAt":          time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Execute update
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		level.Error(r.logger).Log("error", "Failed to update callback status", "id", id, "err", err)
		return err
	}

	level.Info(r.logger).Log("info", "Updated callback status successfully", "id", id, "status", status)
	return nil
}
