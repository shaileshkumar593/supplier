package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/utils/constant"
	"time"

	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PreorderDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	SequenceId      string             `bson:"sequenceId"`
	OtaOrderId      string             `bson:"otaOrderId"`
	Contacts        []trip.Contact     `bson:"contacts"`
	Items           []trip.Item        `bson:"items"`
	RequestCategory string             `bson:"requestCategory"`
	CreatedAt       string             `bson:"createdAt"  validate:"required"`
	UpdatedAt       string             `bson:"updatedAt"`
}

type PreOrderPaymentDocument struct {
	ID                   primitive.ObjectID         `bson:"_id,omitempty"`
	SequenceId           string                     `bson:"sequenceId"`
	OtaOrderId           string                     `bson:"otaOrderId"`
	SupplierOrderId      string                     `bson:"supplierOrderId"`
	ConfirmType          int                        `bson:"confirmType"`
	OrderLastConfirmTime string                     `bson:"orderLastConfirmTime"`
	Items                []trip.PreOrderPaymentItem `bson:"items"`
	Coupons              []trip.PreOrderCoupon      `bson:"coupons"`
	RequestCategory      string                     `bson:"requestCategory"`
	CreatedAt            string                     `bson:"createdAt"  validate:"required"`
	UpdatedAt            string                     `bson:"updatedAt"`
}

type CancellationDocument struct {
	ID              primitive.ObjectID      `bson:"_id,omitempty"`
	SequenceId      string                  `bson:"sequenceId"`
	OTAOrderId      string                  `bson:"otaOrderId"`
	SupplierOrderId string                  `bson:"supplierOrderId"`
	ConfirmType     int                     `bson:"confirmType"`
	Items           []trip.CancellationItem `bson:"items"`
	RequestCategory string                  `bson:"requestCategory"`
	CreatedAt       string                  `bson:"createdAt"  validate:"required"`
	UpdatedAt       string                  `bson:"updatedAt"`
}

// InsertPreorderRequestFromTrip to insert preorder request from trip
func (r *mongoRepository) InsertPreorderRequestFromTrip(ctx context.Context, preorder trip.PreorderRequest) (err error) {
	level.Info(r.logger).Log("repository method", "InsertPreorderRequest")

	collection := r.db.Collection("trip_preorders_request")

	currentTime := time.Now().UTC().Format(time.RFC3339)
	doc := PreorderDocument{
		ID:              primitive.NewObjectID(),
		SequenceId:      preorder.SequenceID,
		OtaOrderId:      preorder.OtaOrderID,
		Contacts:        preorder.Contacts,
		Items:           preorder.Items,
		RequestCategory: constant.TRIPPREORDERREQUEST,
		CreatedAt:       currentTime,
		UpdatedAt:       currentTime,
	}

	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert preorder: %w", err)
	}

	return nil
}

// InsertPaymentRequestFromTrip to insert payment request from trip
func (r *mongoRepository) InsertPaymentRequestFromTrip(ctx context.Context, payment trip.PreOrderPaymentRequest) (err error) {
	level.Info(r.logger).Log("repository method", "InsertPreorderRequest")

	collection := r.db.Collection("trip_payments_request")

	currentTime := time.Now().UTC().Format(time.RFC3339)

	doc := PreOrderPaymentDocument{
		ID:                   primitive.NewObjectID(),
		SequenceId:           payment.SequenceId,
		OtaOrderId:           payment.OtaOrderId,
		SupplierOrderId:      payment.SupplierOrderId,
		ConfirmType:          payment.ConfirmType,
		OrderLastConfirmTime: payment.OrderLastConfirmTime,
		Items:                payment.Items,
		Coupons:              payment.Coupons,
		RequestCategory:      constant.TRIPPAYMENTREQUEST,
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
	}

	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert preorder payment: %w", err)
	}

	return nil
}

// InsertFullCancelOrderRequestFromTrip  to insert full cancel order request
func (r *mongoRepository) InsertFullCancelOrderRequestFromTrip(ctx context.Context, cancelRequest trip.CancellationRequest) (err error) {
	level.Info(r.logger).Log("repository method", "InsertFullCancelOrder")

	collection := r.db.Collection("trip_full_cancel_request")

	currentTime := time.Now().UTC().Format(time.RFC3339)

	doc := CancellationDocument{
		ID:              primitive.NewObjectID(),
		SequenceId:      cancelRequest.SequenceID,
		OTAOrderId:      cancelRequest.OTAOrderID,
		SupplierOrderId: cancelRequest.SupplierOrderID,
		ConfirmType:     cancelRequest.ConfirmType,
		Items:           cancelRequest.Items,
		RequestCategory: constant.TRIPFULLORDERCANCELREQUEST,
		CreatedAt:       currentTime,
		UpdatedAt:       currentTime,
	}

	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert full cancel order: %w", err)
	}

	return nil
}

// FetchTripRequests  to fetch all request data based on collection
/* func (r *mongoRepository) FetchTripRequests(ctx context.Context) (triporders [][]trip.ResponseForTripRequest, err error) {
//level.Info(r.logger).Log("repository method", "FetchTripRequests")

collectionName := []string{"trip_preorders_request", "trip_payments_request", "trip_full_cancel_request"}
triporders = make([][]trip.ResponseForTripRequest, len(collectionName)) // 🔧 Initialize outer slice

DaysAgo := time.Now().AddDate(0, 0, -15) // last 15 days
filter := bson.M{
	"createdAt": bson.M{"$gte": DaysAgo},
}

projection := bson.M{
	"sequenceId":      1,
	"otaOrderId":      1,
	"requestCategory": 1,
	"_id":             0,
}

for i, collName := range collectionName {
	collection := r.db.Collection(collName)
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch documents for coll %s: %w", collName, err)
	}
	defer cursor.Close(ctx)

	var results []trip.ResponseForTripRequest
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents for coll %s: %w", collName, err)
	}
	triporders[i] = results
}
/*
	[][]trip.ResponseForTripRequest{
	{}, // trip_preorders_request
	{}, // trip_payments_request
	{}, // trip_full_cancel_request
*/

/*	return triporders, nil
}
*/

func (r *mongoRepository) FetchTripRequests(ctx context.Context) (triporders [][]trip.ResponseForTripRequest, err error) {
	//level.Info(r.logger).Log("msg", "Fetching trip requests from all collections")

	collections := []string{
		"trip_preorders_request",
		"trip_payments_request",
		"trip_full_cancel_request",
	}

	triporders = make([][]trip.ResponseForTripRequest, len(collections))

	startDate := time.Now().AddDate(0, 0, -15)
	endDate := time.Now()

	// Convert to string for string-based Mongo comparison
	startDateStr := startDate.Format(time.RFC3339)
	endDateStr := endDate.Format(time.RFC3339)

	// Apply string filter instead of time.Time
	filter := bson.M{
		"createdAt": bson.M{
			"$gte": startDateStr,
			"$lte": endDateStr,
		},
	}

	/*level.Info(r.logger).Log(
		"msg", "Applying createdAt filter for string fields",
		"startDate", startDateStr,
		"endDate", endDateStr,
	)*/

	projection := bson.M{
		"sequenceId":      1,
		"otaOrderId":      1,
		"requestCategory": 1,
		"_id":             0,
	}

	for i, collName := range collections {
		//level.Info(r.logger).Log("msg", "Querying MongoDB collection", "collection", collName)

		coll := r.db.Collection(collName)
		cursor, err := coll.Find(ctx, filter, options.Find().SetProjection(projection))
		if err != nil {
			level.Error(r.logger).Log("msg", "MongoDB Find failed", "collection", collName, "error", err.Error())
			return nil, fmt.Errorf("failed to fetch documents for collection %s: %w", collName, err)
		}

		var results []trip.ResponseForTripRequest
		if err := cursor.All(ctx, &results); err != nil {
			cursor.Close(ctx)
			level.Error(r.logger).Log("msg", "MongoDB decode failed", "collection", collName, "error", err.Error())
			return nil, fmt.Errorf("failed to decode documents for collection %s: %w", collName, err)
		}
		cursor.Close(ctx)

		// level.Info(r.logger).Log("msg", "Documents fetched", "collection", collName, "count", len(results))
		triporders[i] = results
	}

	return triporders, nil
}

func (r *mongoRepository) GetSequenceIDByOtaOrderIDAndRequestCategory(ctx context.Context, otaOrderID, serviceName string) (bool, error) {

	filter := bson.M{
		"otaOrderId": otaOrderID,
	}

	var collection *mongo.Collection
	var sequenceID string
	var err error

	switch serviceName {
	case "CreatePreOrder":
		collection = r.db.Collection("trip_preorders_request")
		var doc PreorderDocument
		err = collection.FindOne(ctx, filter).Decode(&doc)
		sequenceID = doc.SequenceId

	case "PayPreOrder":
		collection = r.db.Collection("trip_payments_request")
		var doc PreOrderPaymentDocument
		err = collection.FindOne(ctx, filter).Decode(&doc)
		sequenceID = doc.SequenceId

	case "CancelOrder":
		collection = r.db.Collection("trip_full_cancel_request")
		var doc CancellationDocument
		err = collection.FindOne(ctx, filter).Decode(&doc)
		sequenceID = doc.SequenceId
	}
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil // Document not found is not an error
		}
		r.logger.Log("error", "failed to fetch document", "service", serviceName, "err", err)
		return false, err
	}

	if strings.TrimSpace(sequenceID) == "" {
		r.logger.Log("info", "sequenceId is empty", "service", serviceName, "otaOrderId", otaOrderID)
		return false, nil
	}

	return true, nil
}
