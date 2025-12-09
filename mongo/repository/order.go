package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"swallow-supplier/mongo/domain/yanolja"
	domain "swallow-supplier/mongo/domain/yanolja"
	req_resp "swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/constant"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertPreOrder  insert preorder  record
func (r *mongoRepository) InsertPreOrder(ctx context.Context, rec req_resp.WaitingForOrder) (id string, err error) {
	level.Info(r.logger).Log("repository method ", "InsertPreOrder")

	collection := r.db.Collection("orders")

	// Get current time in UTC and format it as a string
	currentTime := time.Now().UTC().Format(time.RFC3339)

	var selectVariant []domain.SelectVariant
	for _, variant := range rec.SelectVariants {
		// Safely dereference Date and Time if they are not nil
		var date, time string
		if variant.Date != nil {
			date = *variant.Date
		} else {
			date = ""
		}
		if variant.Time != nil {
			time = *variant.Time
		} else {
			time = ""
		}

		record, _ := r.FetchProductByProductId(ctx, variant.ProductID)

		variant := domain.SelectVariant{
			ProductID:        variant.ProductID,
			ProductName:      record.ProductName,
			ProductVersion:   variant.ProductVersion,
			VariantID:        variant.VariantID,
			Date:             date, // assign the string value
			Time:             time, // assign the string value
			Quantity:         variant.Quantity,
			Currency:         variant.Currency,
			PartnerSalePrice: variant.PartnerSalePrice,
			CostPrice:        variant.CostPrice,
		}
		selectVariant = append(selectVariant, variant)
	}
	// JSON data converted to Go struct
	order := domain.Model{
		Id:                            primitive.NewObjectID().Hex(),
		Suppliers:                     strings.ToUpper(constant.SUPPLIERYANOLJA),
		PartnerOrderID:                rec.PartnerOrderID,
		PartnerOrderGroupID:           rec.PartnerOrderGroupID,
		PartnerOrderChannelCode:       rec.PartnerOrderChannelCode,
		PartnerOrderChannelName:       rec.PartnerOrderChannelName,
		TotalSelectedVariantsQuantity: rec.TotalSelectedVariantsQuantity,
		Customer: domain.Customer{
			Name:  rec.Customer.Name,
			Tel:   rec.Customer.Tel,
			Email: rec.Customer.Email,
		},
		ActualCustomer: domain.Customer{
			Name:  rec.ActualCustomer.Name,
			Tel:   rec.ActualCustomer.Tel,
			Email: rec.ActualCustomer.Email,
		},
		SelectVariants: selectVariant,
		OrderExpired:   false,
		OodoSyncStatus: false,
		CreatedAt:      currentTime, // Set CreatedAt timestamp
	}

	// Insert the order into the MongoDB collection
	insertResult, err := collection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(string), nil
}

// UpdatePreOrderById update order by Id
func (r *mongoRepository) UpdatePreOrderById(ctx context.Context, update map[string]any) (id string, err error) {
	level.Info(r.logger).Log("repository method ", "UpdatePreOrderById")

	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}
	if len(update) > 0 {

		for key, val := range update {
			updateBson["$set"].(bson.M)[key] = val
		}
	}

	collection := r.db.Collection("orders")

	filter := bson.M{"_id": update["_id"]}
	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("failed to update document: %w", err)
	}

	val := update["_id"]

	return val.(string), nil
}

// GetOrderbyOrderId retrive single record using orderId
func (r *mongoRepository) GetOrderbyOrderId(ctx context.Context, orderid int64) (record domain.Model, err error) {
	level.Info(r.logger).Log("repository method ", "GetOrderbyOrderId")

	collection := r.db.Collection("orders")

	filter := bson.M{"orderId": orderid}

	err = collection.FindOne(context.TODO(), filter).Decode(&record)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(r.logger).Log("no document exist with orderId: %d", orderid)
			return record, mongo.ErrNoDocuments
		}
		return record, err
	}

	return record, nil
}

// GetOrdersByOdooSyncStatus
func (r *mongoRepository) GetOrdersByOdooSyncStatus(ctx context.Context) ([]domain.Model, error) {

	collection := r.db.Collection("orders")
	// Filter to match only   for
	filter := bson.M{
		"orderStatusCode": "DONE",
		"oodoSyncStatus":  false,
	}

	// Define options to limit result to 100 documents
	findOptions := options.Find()
	findOptions.SetLimit(100) // 100 record at a time

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		level.Error(r.logger).Log("error", "mongo find failed", "details", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []domain.Model
	for cursor.Next(ctx) {
		var record domain.Model
		if err := cursor.Decode(&record); err != nil {
			level.Error(r.logger).Log("error", "mongo decode failed", "details", err.Error())
			return nil, err
		}
		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		level.Error(r.logger).Log("error", "cursor iteration error", "details", err.Error())
		return nil, err
	}

	// If no records found, return a custom error
	if len(records) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return records, nil
}

// GetOrderbyPartnerOrderId
func (r *mongoRepository) GetOrderbyPartnerOrderId(ctx context.Context, partnerOrderId string) (record domain.Model, err error) {
	level.Info(r.logger).Log("repository method ", "GetOrderbyPartnerOrderId")

	collection := r.db.Collection("orders")

	filter := bson.M{"partnerOrderId": partnerOrderId}

	err = collection.FindOne(context.TODO(), filter).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return record, err
		}
	}
	return record, nil
}

// UpdateOrderByOrderId  update single order by orderId
func (r *mongoRepository) UpdateOrderByOrderId(ctx context.Context, orderid int64, update map[string]any) (id string, err error) {
	level.Info(r.logger).Log("repository method ", "UpdateOrderByOrderId")

	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}
	if len(update) > 0 {

		for key, val := range update {
			updateBson["$set"].(bson.M)[key] = val
		}
	}

	collection := r.db.Collection("orders")

	filter := bson.M{"orderId": orderid}

	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("failed to update document: %w", err)
	}

	return string(orderid), nil
}

// UpdateProcessingRestoringOfOrder  for updating orderVariantStatusTypeCode
func (r *mongoRepository) UpdateProcessingRestoringOfOrder(ctx context.Context, orderId int64, updaterec map[string]any) (err error) {
	level.Info(r.logger).Log("repository method ", "UpdateProcessingRestoringOfOrder")

	collection := r.db.Collection("orders")

	// Construct the filter to match orderId, partnerOrderId, and the specific orderVariantId in the orderVariants array
	filter := bson.M{
		"orderId":                      orderId,
		"partnerOrderId":               updaterec["partnerOrderId"],
		"orderVariants.orderVariantId": updaterec["orderVariantId"], // Match specific orderVariantId

	}

	// Create the update document using the positional $ operator to target only the matched orderVariant
	update := bson.M{
		"$set": bson.M{
			"orderVariants.$.orderVariantStatusTypeCode":  updaterec["orderVariantStatusTypeCode"],
			"orderVariants.$.usedRestoreDateTime":         updaterec["dateTime"],
			"orderVariants.$.usedRestoreDateTimeTimezone": updaterec["dateTimeTimeZone"],
			"orderVariants.$.usedRestoreDateTimeOffset":   updaterec["dateTimeOffset"],
			"odooSyncStatus": false,
			"updatedAt":      time.Now().UTC().Format(time.RFC3339), // Add updatedAt timestamp
		},
	}

	// Set options for update
	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	// Execute the UpdateOne to match and update a single orderVariant
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		level.Error(r.logger).Log("repository-error ", "UpdateProcessingRestoringOfOrder")
		return fmt.Errorf("failed to update Resusal to cancel info: %w", err)
	}

	return nil
}

/* func (r *mongoRepository) UpdateProcessingRestoringOfOrder(ctx context.Context, orderId int64, updaterec map[string]any) error {
	level.Info(r.logger).Log("repository method", "UpdateProcessingRestoringOfOrder")

	collection := r.db.Collection("orders")

	// Construct the filter to match orderId and partnerOrderId
	filter := bson.M{
		"orderId":        orderId,
		"partnerOrderId": updaterec["partnerOrderId"],
	}

	// Create the update document to modify all orderVariants
	update := bson.M{
		"$set": bson.M{
			"orderVariants.$[elem].orderVariantStatusTypeCode":  updaterec["orderVariantStatusTypeCode"],
			"orderVariants.$[elem].usedRestoreDateTime":         updaterec["dateTime"],
			"orderVariants.$[elem].usedRestoreDateTimeTimezone": updaterec["dateTimeTimeZone"],
			"orderVariants.$[elem].usedRestoreDateTimeOffset":   updaterec["dateTimeOffset"],
			"updatedAt": time.Now().UTC().Format(time.RFC3339), // Add updatedAt timestamp
		},
	}

	// Define an array filter to match all elements inside orderVariants
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.orderVariantId": bson.M{"$exists": true}}},
	})

	// Execute the UpdateMany to update all orderVariants in matching orders
	_, err := collection.UpdateMany(ctx, filter, update, arrayFilters)
	if err != nil {
		level.Error(r.logger).Log("repository-error", "UpdateProcessingRestoringOfOrder", "error", err)
		return fmt.Errorf("failed to update all orderVariants: %w", err)
	}

	return nil
} */

// DeleteOrderByOrderId  delete the record by orderid
func (r *mongoRepository) DeleteOrderByOrderId(ctx context.Context, orderid int64) (id string, err error) {
	level.Info(r.logger).Log("repository method ", "DeleteOrderByOrderId")

	collection := r.db.Collection("orders")

	filter := bson.M{"orderId": orderid}

	// Use FindOneAndDelete to find and delete the document
	var deletedDoc domain.Model
	err = collection.FindOneAndDelete(ctx, filter).Decode(&deletedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no document found with the given ID %d", orderid)
		}
		return "", fmt.Errorf("failed to delete document with id %d with error %v", orderid, err)
	}

	return deletedDoc.Id, nil
}

// UpdateOrderDueToRefusalToCancel  update single order by orderId,partnerOrderId,orderVariantId
func (r *mongoRepository) UpdateOrderDueToRefusalToCancel(ctx context.Context, orderid int64, partnerOrderId string, orderVariantId int64, update map[string]any) (id string, err error) {
	level.Info(r.logger).Log("repository method ", "UpdateOrderDueToRefusalToCancel")

	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}
	if len(update) > 0 {

		for key, val := range update {
			updateBson["$set"].(bson.M)[key] = val
		}
	}
	collection := r.db.Collection("orders")

	filter := bson.M{"orderId": orderid,
		"partnerOrderId": partnerOrderId,
		"orderVariantId": orderVariantId,
	}

	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("failed to update document: %w", err)
	}

	return string(orderid), nil
}

// GetOrderStatusLookup retrive single record using orderId,partnerorderid,ordervariantid
func (r *mongoRepository) GetOrderStatusLookup(ctx context.Context, orderId int64, partnerOrderId string, orderVariantId int64) ([]domain.OrderVariant, error) {
	level.Info(r.logger).Log("repository method ", "GetOrderStatusLookup")

	collection := r.db.Collection("orders")

	// Define base filter for orderId and partnerOrderId
	filter := bson.M{
		"orderId":        orderId,
		"partnerOrderId": partnerOrderId,
	}

	// If a specific orderVariantId is provided, filter for that variant within the array
	if orderVariantId > 0 {
		filter["orderVariants"] = bson.M{
			"$elemMatch": bson.M{"orderVariantId": orderVariantId},
		}
	}

	// Only project the fields necessary for the orderVariants array to minimize data transfer
	projection := bson.M{
		"orderVariants": 1, // Only retrieve the `orderVariants` array
	}

	// Execute the query
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document exists with orderId: %d", orderId)
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("failed to close MongoDB cursor: %v", cerr)
		}
	}()

	// Parse the results to extract OrderVariant details
	var models []domain.Model
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	// Collect all matched order variants
	var orderVariants []domain.OrderVariant
	for _, model := range models {
		for _, variant := range model.OrderVariants {
			if orderVariantId == 0 || variant.OrderVariantID == orderVariantId {
				orderVariants = append(orderVariants, variant)
			}
		}
	}

	// If no matching variants are found, return an error
	if len(orderVariants) == 0 {
		return nil, fmt.Errorf("no matching order variants found for orderId: %d", orderId)
	}

	return orderVariants, nil
}

// GetOrderStatusLookup retrive single record using orderId,partnerorderid,ordervariantid
func (r *mongoRepository) GetProductIdFromOrder(ctx context.Context, orderid int64, partnerorderid string, ordervariantid int64) (productId int64, err error) {
	level.Info(r.logger).Log("repo-method", "GetProductIdFromOrder")

	collection := r.db.Collection("orders")

	// Construct the filter to match the orderId, partnerOrderId, and find the specific orderVariantId
	filter := bson.M{
		"orderId":                      orderid,
		"partnerOrderId":               partnerorderid,
		"orderVariants.orderVariantId": ordervariantid, // Match directly in the array
	}

	// Set up projection to return only the matched orderVariant
	projection := bson.M{
		"orderVariants": bson.M{"$elemMatch": bson.M{"orderVariantId": ordervariantid}}, // Only the matching element
	}

	// Set options for FindOne including the projection
	findOptions := options.FindOne().SetProjection(projection)

	var orderModel struct {
		OrderVariants []domain.OrderVariant `bson:"orderVariants"`
	}

	err = collection.FindOne(ctx, filter, findOptions).Decode(&orderModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = fmt.Errorf("no document exists with orderId: %d", orderid)
		}
		return -1, err
	}

	// Check if any orderVariant matches the given ordervariantid
	if len(orderModel.OrderVariants) > 0 {
		orderVariant := orderModel.OrderVariants[0] // Get the first matching variant
		level.Info(r.logger).Log("ordervariant received", orderVariant)
		return orderVariant.ProductID, nil
	}

	return -1, fmt.Errorf("no order variant found with orderVariantId: %d", ordervariantid)
}

// UpdateOrderVoucherIndividually  update single order by orderId,partnerOrderId,orderVariantId
func (r *mongoRepository) UpdateOrderVoucherIndividually(ctx context.Context, orderid int64, partnerOrderId string, orderVariantId, orderVariantItemId int64, update map[string]any) (err error) {
	level.Info(r.logger).Log("repo-method", "UpdateOrderVoucherIndividually")

	// Create an empty bson.M map for the update
	updateBson := bson.M{
		"$set": bson.M{
			"oodoSyncStatus": false,
			"updatedAt":      time.Now().Format(time.RFC3339), // Set the updatedAt field to the current time
		},
	}

	// Populate the update map with values from the update parameter
	if len(update) > 0 {
		for key, val := range update {
			updateBson["$set"].(bson.M)[fmt.Sprintf("orderVariants.$[variant].orderVariantItems.$[item].voucher.%s", key)] = val
		}
	}

	// Define the filter to find the correct order and variant
	filter := bson.M{
		"orderId":                      orderid,
		"partnerOrderId":               partnerOrderId,
		"orderVariants.orderVariantId": orderVariantId,
		"orderVariants.orderVariantItems.orderVariantItemId": orderVariantItemId,
	}

	// Set array filters for the update operation
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"variant.orderVariantId": orderVariantId},
			bson.M{"item.orderVariantItemId": orderVariantItemId},
		},
	}

	// Define options for the update
	opts := options.Update().
		SetUpsert(false).
		SetArrayFilters(arrayFilters)

	// Execute the update operation
	result, err := r.db.Collection("orders").UpdateOne(ctx, filter, updateBson, opts)
	if err != nil || result.ModifiedCount == 0 {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// ForcedCancellationReasonUpdate  callback to update cancele status of ordervariantStatusCode
func (r *mongoRepository) ForcedCancellationReasonUpdate(ctx context.Context, orderid int64, partnerOrderId string, orderVariantId int64, ForceCancelReason string) (err error) {

	level.Info(r.logger).Log("repository method", "ForcedCancellationReasonUpdate")

	collection := r.db.Collection("orders")

	// Construct the filter to match orderId, partnerOrderId, and the specific orderVariantId in the orderVariants array
	filter := bson.M{
		"orderId":                      orderid,
		"partnerOrderId":               partnerOrderId,
		"orderVariants.orderVariantId": orderVariantId, // Match specific orderVariantId

	}

	// Create the update document using the positional $ operator to target only the matched orderVariant
	update := bson.M{
		"$set": bson.M{
			"orderVariants.$.forceCancelTypeCode":        ForceCancelReason,
			"orderVariants.$.orderVariantStatusTypeCode": constant.ORDERVARIANTCANCELEDSTATUS,
			"oodoSyncStatus": false,
			"updatedAt":      time.Now().UTC().Format(time.RFC3339), // Add updatedAt timestamp
		},
	}

	// Set options for update
	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	// Execute the UpdateOne to match and update a single orderVariant
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		level.Error(r.logger).Log("repository-error ", "ForcedCancellationReasonUpdate")
		return fmt.Errorf("failed to update Resusal to cancel info: %w", err)
	}

	return nil
}

// GetReconciliationDetailsByOrderAndVariant  retrive reconcilation record by date, order and variant
func (r *mongoRepository) GetReconciliationDetailsByOrderAndVariant(
	ctx context.Context, orderId int64, orderVariantId int64, variantId int64, productId int64,
) ([]yanolja.ReconcilationDetail, error) {

	level.Info(r.logger).Log(
		"repository method ", "GetReconciliationDetailsByOrderAndVariant")

	collection := r.db.Collection("orders")

	// Define the aggregation pipeline stages
	matchStage := bson.D{
		{Key: "$match", Value: bson.M{
			"orderId": orderId,
			"orderVariants": bson.M{
				"$elemMatch": bson.M{
					"orderVariantId": orderVariantId,
					"variantId":      variantId,
					"productId":      productId,
				},
			},
		}},
	}

	// Project only matching variants' `reconciliationByDate` details, grouped by `OrderVariantId`
	projectStage := bson.D{
		{Key: "$project", Value: bson.M{
			"orderVariants": bson.M{
				"$map": bson.M{
					"input": bson.M{
						"$filter": bson.M{
							"input": "$orderVariants",
							"as":    "variant",
							"cond": bson.M{
								"$and": []bson.M{
									{"$eq": bson.A{"$$variant.orderVariantId", orderVariantId}},
									{"$eq": bson.A{"$$variant.variantId", variantId}},
									{"$eq": bson.A{"$$variant.productId", productId}},
								},
							},
						},
					},
					"as": "filteredVariant",
					"in": bson.M{
						"orderVariantId": "$$filteredVariant.orderVariantId",
						"reconciliationByDate": bson.M{
							"$filter": bson.M{
								"input": "$$filteredVariant.reconciliationByDate",
								"as":    "reconciliation",
								"cond": bson.M{
									"$eq": bson.A{"$$reconciliation.reconciliationDate", time.Now().UTC().Format("2006-01-02")},
								},
							},
						},
					},
				},
			},
		}},
	}

	// Run the aggregation
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage}, options.Aggregate().SetBatchSize(1))
	if err != nil {
		return nil, fmt.Errorf("failed to perform aggregation query: %w", err)
	}
	defer cursor.Close(ctx)

	// Structure for decoding
	var allReconciliationDetails []yanolja.ReconcilationDetail
	for cursor.Next(ctx) {
		var result struct {
			OrderVariants []struct {
				OrderVariantId       int64                         `bson:"orderVariantId"`
				ReconciliationByDate []yanolja.ReconcilationDetail `bson:"reconciliationByDate"`
			} `bson:"orderVariants"`
		}

		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode aggregation result: %w", err)
		}

		// Append only relevant reconciliation details
		for _, variant := range result.OrderVariants {
			allReconciliationDetails = append(allReconciliationDetails, variant.ReconciliationByDate...)
		}
	}

	// Check if no relevant reconciliation details were found
	if len(allReconciliationDetails) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return allReconciliationDetails, nil
}

// UpdateCancelDetailsForVariants updates cancel details for multiple order variants.
func (r *mongoRepository) UpdateCancelDetailsForVariants(
	ctx context.Context,
	orderId int64,
	productId int64,
	ordervariantId int64,
	cancelFailReasonCode string,
	cancelStatusCode string,
) error {
	// Specify the collection
	collection := r.db.Collection("orders")

	// Define the filter to locate documents and specific variants within `OrderVariants`
	filter := bson.M{
		"orderId": orderId,
		"orderVariants": bson.M{
			"$elemMatch": bson.M{
				"productId":      productId,
				"orderVariantId": ordervariantId,
			},
		},
	}

	// Define the update operation to set `cancelFailReasonCode` and `cancelStatusCode`
	update := bson.M{
		"$set": bson.M{
			"orderVariants.$[variantFilter].cancelFailReasonCode": cancelFailReasonCode,
			"orderVariants.$[variantFilter].cancelStatusCode":     cancelStatusCode,
		},
	}

	// Define array filters to target the specific `productId` and `variantId`
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{
				"variantFilter.productId":      productId,
				"variantFilter.orderVariantId": ordervariantId,
			},
		},
	}

	// Define options for `UpdateMany` with array filters
	updateOptions := options.Update().SetArrayFilters(arrayFilters)

	// Execute `UpdateMany` to update all matching documents
	res, err := collection.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return fmt.Errorf("failed to update cancel details: %w", err)
	}

	// Log the count of matched and modified documents
	fmt.Printf("Matched %d document(s) and updated %d variant(s)\n", res.MatchedCount, res.ModifiedCount)
	return nil
}

// UpdateReconciliationDetailByDay  update record of reconcilation day by day
func (r *mongoRepository) UpdateReconciliationDetailByDayInsert(ctx context.Context, req map[string]any, updates []domain.ReconcilationDetail) error {
	fmt.Println("Starting UpdateReconciliationDetailByDay")

	// Ensure required fields exist in the request.
	requiredFields := []string{"orderId", "partnerOrderId", "productId", "variantId"}
	for _, field := range requiredFields {
		if _, ok := req[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	fmt.Println("-----------------------1------------------------------\n")
	// Validate the updates slice to ensure it contains valid ReconciliationDetail.
	if len(updates) == 0 {
		return fmt.Errorf("no reconciliation details provided")
	}
	fmt.Println("-----------------------2------------------------------\n")

	// Get the MongoDB collection.
	collection := r.db.Collection("orders")

	// Define the filter to match the document by orderId, partnerOrderId, productId, and variantId within the orderVariants array.
	filter := bson.M{
		"orderId":        req["orderId"],
		"partnerOrderId": req["partnerOrderId"],
		"orderVariants": bson.M{
			"$elemMatch": bson.M{
				"productId": req["productId"],
				"variantId": req["variantId"],
			},
		},
	}
	fmt.Println("-----------------------3------------------------------\n")

	// Ensure the reconciliationByDate array exists for the orderVariant.
	initReconciliationArray := bson.M{
		"$set": bson.M{
			"orderVariants.$[variantFilter].reconciliationByDate": bson.A{},
		},
	}
	arrayFilter := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{
				"variantFilter.productId": req["productId"],
				"variantFilter.variantId": req["variantId"],
			},
		},
	}
	fmt.Println("-----------------------4------------------------------\n")

	_, err := collection.UpdateOne(ctx, filter, initReconciliationArray, options.Update().SetArrayFilters(arrayFilter).SetUpsert(true))
	if err != nil {
		return fmt.Errorf("failed to initialize reconciliationByDate array: %w", err)
	}
	fmt.Println("-----------------------5------------------------------\n")
	fmt.Println("------------------------------------------------------", len(updates), updates)

	// Proceed with updating or inserting reconciliation details.
	for _, detail := range updates {
		reconciliationDate := detail.ReconciliationDate
		if len(reconciliationDate) > 10 {
			reconciliationDate = reconciliationDate[:10] // Extract only the "YYYY-MM-DD" part
		}

		// Parse the string date and normalize to zero time.
		truncatedDate, err := time.Parse("2006-01-02", reconciliationDate)
		if err != nil {
			return fmt.Errorf("invalid date format for reconciliationDate: %w", err)
		}
		truncatedDate = time.Date(truncatedDate.Year(), truncatedDate.Month(), truncatedDate.Day(), 0, 0, 0, 0, truncatedDate.Location())

		// Store only the date (without time) in the reconciliation detail.
		detail.ReconciliationDate = truncatedDate.Format("2006-01-02")

		fmt.Println("-----------------------6------------------------------\n")

		// Update existing reconciliation record for the same date (ignoring time).
		updateQuery := bson.M{
			"$set": bson.M{
				"orderVariants.$[variantFilter].reconciliationByDate.$[reconFilter].reconcileOrderStatusCode": detail.ReconcileOrderStatusCode,
			},
		}
		fmt.Println("-----------------------7------------------------------\n")

		// Insert new reconciliation detail if it doesn't exist for the date.
		insertQuery := bson.M{
			"$push": bson.M{
				"orderVariants.$[variantFilter].reconciliationByDate": detail,
			},
		}
		fmt.Println("-----------------------8------------------------------\n")

		// Define the array filters for matching the correct variant and reconciliation date (by day only).
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{
					"variantFilter.productId": req["productId"],
					"variantFilter.variantId": req["variantId"],
				},
				bson.M{
					"reconFilter.reconciliationDate": truncatedDate,
				},
			},
		}
		fmt.Println("-----------------------9------------------------------\n")

		// Attempt to update the reconciliation detail if it exists.
		opts := options.Update().SetArrayFilters(arrayFilters).SetUpsert(false)
		fmt.Printf("Attempting to update reconciliation for date %s\n", truncatedDate)

		res, err := collection.UpdateOne(ctx, filter, updateQuery, opts)
		if err != nil {
			return fmt.Errorf("failed to update reconciliation details: %w", err)
		}
		fmt.Println("-----------------------10------------------------------\n")

		// If no document was updated, insert a new reconciliation detail.
		if res.ModifiedCount == 0 {
			fmt.Printf("No existing record for date %s, inserting new record\n", truncatedDate)
			_, err = collection.UpdateOne(ctx, filter, insertQuery, options.Update().SetArrayFilters(arrayFilter).SetUpsert(true))
			if err != nil {
				return fmt.Errorf("failed to insert new reconciliation detail: %w", err)
			}
		}
		fmt.Println("-----------------------11------------------------------\n")

	}
	fmt.Println("-----------------------1222222------------------------------\n")

	return nil
}

// UpdateReconciliationDetailByDay  update record of reconcilation day by day
func (r *mongoRepository) UpdateReconciliationDetailByDay(ctx context.Context, req map[string]any, updates []domain.ReconcilationDetail) error {
	level.Info(r.logger).Log("repository_method_name ", "UpdateReconciliationDetailByDay")

	// Ensure required fields exist in the request.
	requiredFields := []string{"orderId", "partnerOrderId", "productId", "variantId"}
	for _, field := range requiredFields {
		if _, ok := req[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate the updates slice to ensure it contains valid ReconciliationDetail.
	if len(updates) == 0 {
		return fmt.Errorf("no reconciliation details provided")
	}

	// Get the MongoDB collection.
	collection := r.db.Collection("orders")

	// Define the filter to match the document by orderId, partnerOrderId, productId, and variantId within the orderVariants array.
	filter := bson.M{
		"orderId":        req["orderId"],
		"partnerOrderId": req["partnerOrderId"],
		"orderVariants": bson.M{
			"$elemMatch": bson.M{
				"productId": req["productId"],
				"variantId": req["variantId"],
			},
		},
	}

	// Iterate through updates and process each reconciliation detail
	for _, detail := range updates {
		// Set only the date part of reconciliationDate
		reconciliationDate := detail.ReconciliationDate
		if len(reconciliationDate) > 10 {
			reconciliationDate = reconciliationDate[:10] // Extract "YYYY-MM-DD" part only
		}
		detail.ReconciliationDate = reconciliationDate
		// First, attempt to update an existing reconciliation entry with the same reconciliationDate.
		updateQuery := bson.M{
			"$set": bson.M{
				"orderVariants.$[variantFilter].reconciliationByDate.$[reconFilter].reconcileOrderStatusCode": detail.ReconcileOrderStatusCode,
				"updatedAt": time.Now().UTC().Format(time.RFC3339),
			},
		}

		// Define array filters for updating the reconciliation record with the specified date.
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{
					"variantFilter.productId": req["productId"],
					"variantFilter.variantId": req["variantId"],
				},
				bson.M{
					"reconFilter.reconciliationDate": reconciliationDate,
				},
			},
		}

		// Try updating the existing entry first.
		opts := options.Update().SetArrayFilters(arrayFilters).SetUpsert(false)
		res, err := collection.UpdateOne(ctx, filter, updateQuery, opts)
		if err != nil {
			return fmt.Errorf("failed to update reconciliation details: %w", err)
		}

		// If no document was updated (i.e., no matching date entry was found), proceed to insert a new entry.
		if res.ModifiedCount == 0 {
			insertQuery := bson.M{
				"$push": bson.M{
					"orderVariants.$[variantFilter].reconciliationByDate": bson.M{
						"reconciliationDate":       reconciliationDate,
						"reconcileOrderStatusCode": detail.ReconcileOrderStatusCode,
					},
				},
			}

			// Use array filters for pushing a new entry to the specified variant.
			insertArrayFilters := options.ArrayFilters{
				Filters: []interface{}{
					bson.M{
						"variantFilter.productId": req["productId"],
						"variantFilter.variantId": req["variantId"],
					},
				},
			}

			_, err = collection.UpdateOne(ctx, filter, insertQuery, options.Update().SetArrayFilters(insertArrayFilters).SetUpsert(true))
			if err != nil {
				return fmt.Errorf("failed to insert new reconciliation detail: %w", err)
			}
		}
	}

	return nil
}

// GetTotalQuantityByPerson aggregates the total quantity based on name, email,tel,
func (r *mongoRepository) GetTotalQuantityPurchasedByPersonToday(
	ctx context.Context,
	name, email, tel string,
	productId, variantId int64,
) (int32, error) {
	level.Info(r.logger).Log("repository method", "GetTotalQuantityPurchasedByPersonToday")

	today := time.Now().UTC().Truncate(24 * time.Hour)

	pipeline := mongo.Pipeline{
		// Stage 1: Convert createdAt to date and truncate to day
		{{
			Key: "$addFields",
			Value: bson.D{
				{"createdDay", bson.D{
					{"$dateTrunc", bson.D{
						{"date", bson.D{
							{"$dateFromString", bson.D{
								{"dateString", "$createdAt"},
							}},
						}},
						{"unit", "day"},
					}},
				}},
			},
		}},
		// Stage 2: Match customer, product, and date
		{{
			Key: "$match",
			Value: bson.D{
				{"$expr", bson.D{
					{"$eq", bson.A{bson.D{{"$toLower", "$customer.name"}}, strings.ToLower(name)}},
				}},
				{"$expr", bson.D{
					{"$eq", bson.A{bson.D{{"$toLower", "$customer.email"}}, strings.ToLower(email)}},
				}},
				{"customer.tel", tel},
				{"createdDay", today},
				{"orderStatusCode", "DONE"},
			},
		}},
		// Stage 3: Unwind selectVariants
		{{Key: "$unwind", Value: "$selectVariants"}},
		// Stage 4: Match product/variant
		{{
			Key: "$match",
			Value: bson.D{
				{"selectVariants.productId", productId},
				{"selectVariants.variantId", variantId},
			},
		}},
		// Stage 5: Group
		{{
			Key: "$group",
			Value: bson.D{
				{"_id", nil},
				{"totalQuantity", bson.D{{"$sum", "$selectVariants.quantity"}}},
			},
		}},
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("aggregation error: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			TotalQuantity int32 `bson:"totalQuantity"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("decode error: %w", err)
		}
		return result.TotalQuantity, nil
	}

	return 0, nil
}

/*func (r *mongoRepository) GetTotalQuantityPurchasedByPerson(ctx context.Context, name, email, tel string, productId int64, variantId int64) (totalQuantity int32, err error) {
	level.Info(r.logger).Log("repository method ", "GetTotalQuantityPurchasedByPerson")

	// Create an optimized pipeline for aggregation
	pipeline := mongo.Pipeline{
		// Stage 1: Match documents based on customer and product/variant information
		{
			{Key: "$match", Value: bson.D{
				{"$expr", bson.D{
					{"$eq", bson.A{
						bson.D{{"$toLower", "$customer.name"}}, strings.ToLower(name), // Case-insensitive match for name
					}},
				}},
				{"$expr", bson.D{
					{"$eq", bson.A{
						bson.D{{"$toLower", "$customer.email"}}, strings.ToLower(email), // Case-insensitive match for email
					}},
				}},
				{Key: "customer.tel", Value: tel},       // Matches the customer's telephone
				{"selectVariants.productId", productId}, // Matches the ProductId in SelectVariants
				{"selectVariants.variantId", variantId}, // Matches the VariantId in SelectVariants
			}},
		},
		// Stage 2: Unwind the selectVariants array (deconstruct it)
		{
			{"$unwind", "$selectVariants"},
		},
		// Stage 3: Group the documents by summing up the quantities from selectVariants
		{
			{"$group", bson.D{
				{"_id", nil}, // Grouping all documents together
				{"totalQuantity", bson.D{{"$sum", "$selectVariants.quantity"}}}, // Summing the quantities
			}},
		},
		// Stage 4: Optionally project only the necessary fields
		{
			{"$project", bson.D{
				{"totalQuantity", 1}, // Output the totalQuantity
			}},
		},
	}

	// Execute the aggregation query
	cursor, err := r.db.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("error aggregating total quantity: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the result
	if cursor.Next(ctx) {
		var result struct {
			TotalQuantity int32 `bson:"totalQuantity"` // Total quantity
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("error decoding result: %w", err)
		}
		totalQuantity = result.TotalQuantity
	}

	// Return the totalQuantity or error if not found
	if totalQuantity == 0 {
		return 0, fmt.Errorf("no total quantity found for customer %s with productId %d and variantId %d", name, productId, variantId)
	}

	return totalQuantity, nil
}*/

// GetReconciliationDetailsByDateAndStatus  retrive reconcilation status day wise
func (r *mongoRepository) GetReconciliationDetailsByDateAndStatus(ctx context.Context, reconciliationDate string, statusCode string) (results []req_resp.OrderReconcilation, err error) {
	level.Info(r.logger).Log("repository method", "GetReconciliationDetailsByDateAndStatus")

	// Match stage to filter reconciliation details by date and status
	matchStage := bson.D{
		{"$match", bson.M{
			"orderVariants.reconciliationByDate": bson.M{
				"$elemMatch": bson.M{
					"reconciliationDate":       reconciliationDate,
					"reconcileOrderStatusCode": statusCode,
				},
			},
		}},
	}

	// Projection stage to include `selectVariants` and map required fields
	projectStage := bson.D{
		{"$project", bson.M{
			"orderId":                 1,
			"partnerOrderId":          1,
			"partnerOrderChannelPin":  1,
			"partnerOrderChannelName": 1,
			"partnerOrderChannelCode": 1,
			"selectVariants": bson.M{
				"$map": bson.M{
					"input": "$selectVariants",
					"as":    "variant",
					"in": bson.M{
						"partnerSalePrice": "$$variant.partnerSalePrice", // Map `partnerSalePrice` from `selectVariants`
					},
				},
			},
			"orderVariants": bson.M{
				"$map": bson.M{
					"input": "$orderVariants",
					"as":    "variant",
					"in": bson.M{
						"productId":      "$$variant.productId",
						"variantId":      "$$variant.variantId",
						"orderVariantId": "$$variant.orderVariantId",
						"reconciliationByDate": bson.M{
							"$filter": bson.M{
								"input": "$$variant.reconciliationByDate",
								"as":    "reconDetail",
								"cond": bson.M{
									"$and": []bson.M{
										{"$eq": bson.A{"$$reconDetail.reconciliationDate", reconciliationDate}},
										{"$eq": bson.A{"$$reconDetail.reconcileOrderStatusCode", statusCode}},
									},
								},
							},
						},
					},
				},
			},
		}},
	}

	// Execute aggregation pipeline
	cursor, err := r.db.Collection("orders").Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		return nil, fmt.Errorf("error performing aggregation query: %w", err)
	}
	defer cursor.Close(ctx)

	// Parse results
	for cursor.Next(ctx) {
		var order struct {
			OrderId                 int64  `bson:"orderId"`
			PartnerOrderId          string `bson:"partnerOrderId"`
			PartnerOrderChannelPin  string `bson:"partnerOrderChannelPin"`
			PartnerOrderChannelName string `bson:"partnerOrderChannelName"`
			PartnerOrderChannelCode string `bson:"partnerOrderChannelCode"`
			SelectVariants          []struct {
				PartnerSalePrice float32 `bson:"partnerSalePrice"`
			} `bson:"selectVariants"`
			OrderVariants []struct {
				ProductId            int64                         `bson:"productId"`
				VariantId            int64                         `bson:"variantId"`
				OrderVariantId       int64                         `bson:"orderVariantId"`
				ReconciliationByDate []yanolja.ReconcilationDetail `bson:"reconciliationByDate"`
			} `bson:"orderVariants"`
		}

		if err := cursor.Decode(&order); err != nil {
			return nil, fmt.Errorf("error decoding order document: %w", err)
		}

		for _, variant := range order.OrderVariants {
			orderId := fmt.Sprintf("%d", order.OrderId)
			orderVariantId := fmt.Sprintf("%d", variant.OrderVariantId)
			for _, reconDetail := range variant.ReconciliationByDate {
				// Retrieve partnerSalePrice from `selectVariants`
				var partnerSalePrice float32
				if len(order.SelectVariants) > 0 {
					partnerSalePrice = order.SelectVariants[0].PartnerSalePrice
				}

				result := req_resp.OrderReconcilation{
					ReconciliationDate:         reconDetail.ReconciliationDate,
					ProductId:                  variant.ProductId,
					VariantId:                  variant.VariantId,
					OrderVariantId:             orderVariantId,
					PartnerSalePrice:           partnerSalePrice, // Use `SelectVariants.PartnerSalePrice`
					ReconcileOrderStatusCode:   reconDetail.ReconcileOrderStatusCode,
					PartnerOrderId:             order.PartnerOrderId,
					Pin:                        orderId + orderVariantId,
					PartnerOrderChannelPin:     orderVariantId,
					PartnerOrderChannelName:    order.PartnerOrderChannelName,
					PartnerOrderChannelCode:    order.PartnerOrderChannelCode,
					PartnerOrderChannelOrderId: orderId,
				}
				results = append(results, result)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through aggregation cursor: %w", err)
	}

	return results, nil
}

// UpdateOrderDueToRefusalToCancel  update cancel status of single order
func (r *mongoRepository) UpdateForcedCancelOrderDetail(ctx context.Context, partnerOrderId string, newStatus string) (id string, err error) {
	level.Info(r.logger).Log("repository method ", "UpdateForcedCancelOrderDetail")

	collection := r.db.Collection("orders")

	// Define the filter to match the document with the specified partnerOrderId
	filter := bson.M{
		"partnerOrderId": partnerOrderId,
	}

	// Create the update document to set orderVariantStatusTypeCode for all items in the OrderVariants array
	update := bson.M{
		"$set": bson.M{
			"orderVariants.$[].orderVariantStatusTypeCode": newStatus,
			"oodoSyncStatus": false,
			"updatedAt":      time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Set options for update
	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	// Execute the UpdateOne to match and update the document
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil || result.ModifiedCount == 0 {
		level.Error(r.logger).Log("repository-error", "UpdateForcedCancelOrderDetail", "error", err)
		return "", fmt.Errorf("failed to update document: %w", err)
	}

	return partnerOrderId, nil
}

// UpdateOrderVariantStatusByOrderId  to update all variant status
func (r *mongoRepository) UpdateOrderVariantStatusByOrderId(ctx context.Context, orderid int64, variantStatus string) (err error) {
	level.Info(r.logger).Log("repository method ", "UpdateOrderVariantStatusByOrderId")

	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}

	collection := r.db.Collection("orders")

	filter := bson.M{"orderId": orderid}

	updateBson = bson.M{
		"$set": bson.M{
			"orderVariants.$[].orderVariantStatusTypeCode": variantStatus,
			"oodoSyncStatus": false,
			"updatedAt":      time.Now().UTC().Format(time.RFC3339),
		},
	}

	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil || result.ModifiedCount == 0 {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// UpdateOrderCancelTypeCode  to update all variant canceltypecode
func (r *mongoRepository) UpdateOrderCancelAck(ctx context.Context, orderid int64, partnerOrderId string, canceltypecode string, orderstatus string) (err error) {
	level.Info(r.logger).Log("repository method ", "UpdateOrderCancelAck")

	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}

	collection := r.db.Collection("orders")

	filter := bson.M{"orderId": orderid,
		"partnerOrderId": partnerOrderId}

	updateBson = bson.M{
		"$set": bson.M{
			"orderVariants.$[].orderCancelTypeCode":                                    canceltypecode,
			"orderVariants.$[].canceledDateTime":                                       time.Now().UTC().Format(time.RFC3339),
			"orderVariants.$[].refundInfo":                                             canceltypecode,
			"orderVariants.$[].orderVariantStatusTypeCode":                             orderstatus,
			"orderVariants.$[].orderVariantItems.$[].voucher.voucherProvideStatusCode": orderstatus,
			"oodoSyncStatus": false,
			"updatedAt":      time.Now().UTC().Format(time.RFC3339),
		},
	}

	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	result, err := collection.UpdateMany(context.TODO(), filter, updateBson, opts)
	if err != nil || result.ModifiedCount == 0 {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// UpdateRefusalToCancelInfo
func (r *mongoRepository) UpdateRefusalToCancelInfo(ctx context.Context, orderid int64, partnerOrderId string, ovariantId int64, cancelRejectTypeCode, message string) error {
	level.Info(r.logger).Log("repo-method", "UpdateRefusalToCancelInfo")

	collection := r.db.Collection("orders")

	// Construct the filter for the specific document and orderVariantId
	filter := bson.M{
		"orderId":        orderid,
		"partnerOrderId": partnerOrderId,
		"orderVariants": bson.M{
			"$elemMatch": bson.M{
				"orderVariantId": ovariantId,
			},
		},
	}

	// Log the filter for debugging purposes
	level.Info(r.logger).Log("filter", filter)

	// Construct the update query using the positional operator ($)
	update := bson.M{
		"$set": bson.M{
			"orderVariants.$.cancelRejectTypeCode":       cancelRejectTypeCode,
			"orderVariants.$.message":                    message,
			"orderVariants.$.orderVariantStatusTypeCode": constant.ORDERVARIANTNOTUSEDSTATUS,
			"oodoSyncStatus":                             false,
			"updatedAt":                                  time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Log the update query
	level.Info(r.logger).Log("update", update)

	// Execute the update operation
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		level.Error(r.logger).Log("repository-error", "UpdateRefusalToCancelInfo", "error", err)
		return fmt.Errorf("failed to update refusal to cancel info: %w", err)
	}

	// Log the update result
	level.Info(r.logger).Log("update-result", result.ModifiedCount)

	// Handle no matched document case
	if result.ModifiedCount == 0 {
		level.Warn(r.logger).Log("warning", "No document updated. Verify the filter matches the target document.")
		return fmt.Errorf("no document matched the filter criteria")
	}

	return nil
}

// UpdateForceCancelVariants updates the forceCancelTypeCode for specific order variants in an order.
func (r *mongoRepository) UpdateForceCancelVariants(ctx context.Context, orderid int64, partnerOrderId string, forceCancelVariants []req_resp.CancelledVariants) (err error) {
	level.Info(r.logger).Log("repo-method", "UpdateForceCancelVariants")

	collection := r.db.Collection("orders")

	// Construct array filters for each specific orderVariantId in the orderVariants array.
	arrayFilters := bson.A{}
	for _, variant := range forceCancelVariants {
		filter := bson.M{
			"elem.orderVariantId": variant.OrderVariantID,
		}
		arrayFilters = append(arrayFilters, filter)
	}

	// Construct the update statement to set forceCancelTypeCode for each matched orderVariant.
	update := bson.M{
		"$set": bson.M{
			"orderVariants.$[elem].forceCancelTypeCode": bson.M{
				"$each": forceCancelVariants,
			},
			"oodoSyncStatus": false,
			"updatedAt":      time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Options to apply array filters and prevent upsert.
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: arrayFilters,
	}).SetUpsert(false)

	// Construct the filter for the main document using orderId and partnerOrderId.
	filter := bson.M{
		"orderId":        orderid,
		"partnerOrderId": partnerOrderId,
	}

	// Execute the UpdateMany operation.
	_, err = collection.UpdateMany(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update force cancel variants reason: %w", err)
	}

	return nil
}

func (r *mongoRepository) FindOrderByOrderIdAndPartnerOrderId(ctx context.Context, orederId int64, partnerOrderId string) (record yanolja.Model, err error) {
	level.Info(r.logger).Log("repository method ", "FindOrderByOrderIdAndPartnerOrderId")

	collection := r.db.Collection("orders")

	filter := bson.M{
		"orderId":        orederId,
		"partnerOrderId": partnerOrderId,
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&record)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = fmt.Errorf("no document exist with orderId: %d and partnerorderId %s", orederId, partnerOrderId)
			return record, err
		}
	}

	return record, nil
}

// GetOrdersByChannelCodeAndCustomerEmail
func (r *mongoRepository) GetOrdersByChannelCodeAndCustomerEmail(ctx context.Context, channelCode string, customerEmail string) ([]domain.Model, error) {
	level.Info(r.logger).Log("repository method", "GetOrdersByChannelCodeAndCustomerEmail")

	fmt.Println(" channelCode :", channelCode, "cutomerEmail : ", customerEmail)
	collection := r.db.Collection("orders")

	filter := bson.M{
		"partnerOrderChannelCode": channelCode,
		"customer.email":          customerEmail,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []domain.Model
	for cursor.Next(ctx) {
		var record domain.Model
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// If no records found, return a custom error
	if len(records) == 0 {
		level.Error(r.logger).Log("error ", "no documents(order) exist with odooSyncStatus=false from method GetOrdersByOdooSyncStatus")
		return nil, mongo.ErrNoDocuments
	}

	return records, nil
}

// GetOrderByPartnerIdSuffix
func (r *mongoRepository) GetOrderByPartnerIdSuffix(ctx context.Context, suffix string) (domain.Model, error) {
	level.Info(r.logger).Log("repository method", "GetOrderByPartnerIdSuffix")

	collection := r.db.Collection("orders")

	// Match partnerOrderId ending in the unique suffix
	filter := bson.M{
		"partnerOrderId": bson.M{"$regex": suffix + "$"},
	}

	var record domain.Model
	err := collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Warn(r.logger).Log("warn", "no order found with partnerOrderId suffix", "suffix", suffix)
			return record, err
		}
		level.Error(r.logger).Log("error", "failed to find order with partnerOrderId suffix", "details", err.Error())
		return record, err
	}

	return record, nil
}
