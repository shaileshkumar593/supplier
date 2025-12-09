package repository

import (
	"context"
	"fmt"
	"swallow-supplier/utils/constant"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *mongoRepository) GetSequenceIDByKey(ctx context.Context, key string, typeOfCollection string) (string, error) {
	level.Info(r.logger).Log("repository method", "GetSequenceIDByKey")
	level.Info(r.logger).Log("key", key)

	/* // Remove suffix "-Payment"
	otaOrderId := strings.TrimSuffix(key, "-Payment")
	level.Info(r.logger).Log("otaOrderId", otaOrderId) */

	// Set a safe query timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var collection *mongo.Collection
	if typeOfCollection == constant.TRIPPAYMENTREQUEST {
		collection = r.db.Collection("trip_payments_request")

	} else if typeOfCollection == constant.TRIPFULLORDERCANCELREQUEST {
		collection = r.db.Collection("trip_full_cancel_request")
	}

	// MongoDB query
	filter := bson.M{"otaOrderId": key}

	// Target result struct
	var result struct {
		SequenceID string `bson:"sequenceId"`
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		level.Error(r.logger).Log("error", fmt.Sprintf("document not exist with otaOrderId %s", key))
		return "", err
	}
	r.logger.Log("sequenceId from repository", result.SequenceID)
	return result.SequenceID, nil
}
