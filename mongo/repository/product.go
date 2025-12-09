package repository

import (
	"context"
	"fmt"
	"swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	req_resp "swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/constant"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertProducts
func (r *mongoRepository) InsertProducts(ctx context.Context, products []yanolja.Product) (err error) {
	level.Info(r.logger).Log(
		"operation", "InsertProducts",
	)

	for i := 0; i < len(products); i++ {
		products[i].Id = primitive.NewObjectID().Hex()
	}
	docs := make([]interface{}, len(products))
	for i, prodts := range products {
		prodts.SupplierName = "Yanolja"
		docs[i] = prodts
	}

	collection := r.db.Collection("products")
	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

// InsertOneProduct
func (r *mongoRepository) InsertOneProduct(ctx context.Context, product yanolja.Product) (err error) {
	level.Info(r.logger).Log(
		"operation", "InsertOneProduct",
		"productId", product.ProductID,
	)
	product.Id = primitive.NewObjectID().Hex()
	product.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	product.SupplierName = constant.SUPPLIERYANOLJA

	collection := r.db.Collection("products")
	_, err = collection.InsertOne(ctx, product)
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return err
	}
	return nil
}

// UpsertProduct
func (r *mongoRepository) UpsertProduct(ctx context.Context, product yanolja.Product) error {
	level.Info(r.logger).Log(
		"operation", "UpsertProduct",
		"productId", product.ProductID,
	)

	rsp := make(map[string]int64)
	// Define the current timestamp
	currentTime := time.Now().UTC().Format(time.RFC3339)

	// Set `UpdatedAt` for all operations
	product.UpdatedAt = currentTime

	// If inserting, set `CreatedAt`
	if product.Id == "" {
		product.Id = primitive.NewObjectID().Hex()
		product.CreatedAt = currentTime
	}

	// Define the filter to find the product by ProductID
	filter := bson.M{"productId": product.ProductID}

	// Define the update document
	update := bson.M{
		"$set": bson.M{
			"supplierName":                constant.SUPPLIERYANOLJA,
			"productId":                   product.ProductID,
			"productName":                 product.ProductName,
			"productVersion":              product.ProductVersion,
			"price":                       product.Price,
			"productStatusCode":           product.ProductStatusCode,
			"productTypeCode":             product.ProductTypeCode,
			"salePeriod":                  product.SalePeriod,
			"productBriefIntroduction":    product.ProductBriefIntroduction,
			"productInfo":                 product.ProductInfo,
			"productOptionGroups":         product.ProductOptionGroups,
			"searchKeywords":              product.SearchKeywords,
			"categories":                  product.Categories,
			"regions":                     product.Regions,
			"images":                      product.Images,
			"textFromImages":              product.TextFromImages,
			"videos":                      product.Videos,
			"pictograms":                  product.Pictograms,
			"isCancelPenalty":             product.IsCancelPenalty,
			"isReservationAfterPurchase":  product.IsReservationAfterPurchase,
			"purchaseDateUsableTypeCode":  product.PurchaseDateUsableTypeCode,
			"isAvailableOnPurchaseDate":   product.IsAvailableOnPurchaseDate,
			"isIntegratedVoucher":         product.IsIntegratedVoucher,
			"isRefundableAfterExpiration": product.IsRefundableAfterExpiration,
			"isUsed":                      product.IsUsed,
			"sellerInfos":                 product.SellerInfos,
			"convenienceTypeCode":         product.ConvenienceTypeCode,
			"imageScheduleStatus":         false,
			"viewScheduleStatus":          false,
			"contentScheduleStatus":       false,
			"oodoSyncStatus":              false,
			"updatedAt":                   product.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"createdAt": product.CreatedAt,
		},
	}

	// Specify the upsert option (insert if not found)
	opts := options.Update().SetUpsert(true)

	// Perform the update operation
	collection := r.db.Collection("products")
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		level.Error(r.logger).Log("database error", err.Error())
		return fmt.Errorf("failed to upsert product: %w", err)
	}
	level.Info(r.logger).Log(
		"updated_count", result.ModifiedCount,
		"inserted_count", result.UpsertedCount,
		"matched_count", result.MatchedCount,
	)

	rsp["updated_count"] = result.ModifiedCount
	rsp["inserted_count"] = result.UpsertedCount
	rsp["matched_count"] = result.MatchedCount

	//fmt.Println("************ Detail***************** : ", rsp)
	return nil
}

// UpdateProductsByProductID
func (r *mongoRepository) UpdateProductsByProductID(ctx context.Context, productId int64, update map[string]any) (id string, err error) {
	level.Info(r.logger).Log(
		"operation", "UpdateProductsByProductID",
		"productId", productId,
	)
	// Create an empty bson.M map
	updateBson := bson.M{
		"$set": bson.M{},
	}
	if len(update) > 0 {

		for key, val := range update {
			updateBson["$set"].(bson.M)[key] = val
		}
	}

	collection := r.db.Collection("products")

	filter := bson.M{"productId": productId}

	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found
	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil {
		return "", fmt.Errorf("failed to update document: %w", err)
	}

	return result.UpsertedID.(string), nil
}

// FetchProductByProductId
func (r *mongoRepository) FetchProductByProductId(ctx context.Context, productId int64) (record yanolja.Product, err error) {

	level.Info(r.logger).Log(
		"operation", "FetchProductByProductId",
		"productId", productId,
	)
	filter := bson.M{"productId": productId}

	collection := r.db.Collection("products")

	err = collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return record, mongo.ErrNoDocuments
		}
		return record, err
	}
	return record, nil
}

// FetchProducts
func (r *mongoRepository) FetchProducts(ctx context.Context) (records []yanolja.Product, err error) {
	level.Info(r.logger).Log(
		"operation", "FetchProducts",
	)
	filter := bson.M{} // Empty filter to match all documents

	collection := r.db.Collection("products")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product yanolja.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		records = append(records, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// FetchProductsBasedOnOdooSyncStatus  fetch record based on odooSyncStatus false
func (r *mongoRepository) FetchProductsBasedOnOdooSyncStatus(ctx context.Context) ([]yanolja.Product, error) {
	level.Info(r.logger).Log(
		"operation", "FetchProductsBasedOnOdooSyncStatus",
	)
	filter := bson.M{"oodoSyncStatus": false} // Filter products with odooSyncStatus=false

	collection := r.db.Collection("products")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []yanolja.Product
	for cursor.Next(ctx) {
		var product yanolja.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		records = append(records, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	// If no records found, return a custom error
	if len(records) == 0 {
		return nil, fmt.Errorf("no documents(products) exist with odooSyncStatus=false")
	}

	return records, nil
}

// DeleteProductByProductId
func (r *mongoRepository) DeleteProductByProductId(ctx context.Context, productId int64) (id string, err error) {
	level.Info(r.logger).Log(
		"operation", "DeleteProductByProductId",
	)

	collection := r.db.Collection("products")

	filter := bson.M{"productId": productId}

	// Use FindOneAndDelete to find and delete the document
	var deletedDoc yanolja.Product
	err = collection.FindOneAndDelete(ctx, filter).Decode(&deletedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no document found with the given ID %d", productId)
		}
		return "", fmt.Errorf("failed to delete document with id %d with error %v", productId, err)
	}

	return deletedDoc.Id, nil
}

// fetch product updated today (inventory updated )
func (r *mongoRepository) FetchAllProductsWithinDateRange(ctx context.Context) (products []yanolja.Product, err error) {

	level.Info(r.logger).Log(
		"operation", "FetchAllProductsWithinDateRange",
	)

	collection := r.db.Collection("products")

	// Get today's date in UTC and calculate the range
	now := time.Now().UTC()
	startOfDay := now.AddDate(0, 0, -20) // Start 20 days ago
	endOfDay := now                      // End date is the current time

	// Convert time.Time to ISO 8601 string format
	startDateStr := startOfDay.Format(time.RFC3339)
	endDateStr := endOfDay.Format(time.RFC3339)

	// Query filter
	filter := bson.M{
		"$or": []bson.M{
			{
				"updatedAt": bson.M{
					"$gte": startDateStr,
					"$lt":  endDateStr,
				},
			},
			{
				"$and": []bson.M{
					{"updatedAt": bson.M{"$eq": ""}},
					{"createdAt": bson.M{
						"$gte": startDateStr,
						"$lt":  endDateStr,
					}},
				},
			},
		},
	}

	// Perform the query
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results into the products slice
	for cursor.Next(ctx) {
		var product yanolja.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, product)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return products, nil
}

// FetchProductImagesForLast12Hours
func (r *mongoRepository) FetchProductImagesForLast12Hours(ctx context.Context) (results []req_resp.ProductImages, err error) {
	level.Info(r.logger).Log(
		"operation", "FetchProductImagesForLast12Hours",
	)

	collection := r.db.Collection("products")

	// Get current time in UTC (end of range)
	endOfRange := time.Now().UTC()

	// Start time: 12 hours before current time
	startOfRange := endOfRange.Add(-12 * time.Hour)

	// Convert times to ISO 8601 string format
	startDateStr := startOfRange.Format(time.RFC3339)
	endDateStr := endOfRange.Format(time.RFC3339)

	// Query filter
	filter := bson.M{
		"$or": []bson.M{
			{"updatedAt": bson.M{"$gte": startDateStr, "$lte": endDateStr}},
			{"createdAt": bson.M{"$gte": startDateStr, "$lte": endDateStr}},
		},
	}

	// Define projection to fetch only ProductID and Images
	projection := bson.M{
		"productId": 1,
		"images":    1,
	}

	// Perform the query with projection
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode each product into result structure
	for cursor.Next(ctx) {
		var result struct {
			ProductID int64           `json:"productId"`
			Images    []yanolja.Image `json:"images"`
		}

		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		results = append(results, result)
	}
	fmt.Println("")
	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}

// FetchProductsUpdatedToday retrieves all products whose createdAt or updatedAt timestamps fall within today's range.
func (r *mongoRepository) FetchProductsWithContentScheduleStatusFalse(ctx context.Context) ([]yanolja.Product, error) {
	level.Info(r.logger).Log("repository", "FetchProductsWithContentScheduleStatusFalse")

	collection := r.db.Collection("products")

	// Enforce timeout to avoid long-running queries
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Filter: only products with contentScheduleStatus == false
	filter := bson.M{"contentScheduleStatus": false}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		level.Error(r.logger).Log("repository", "FetchProductsWithContentScheduleStatusFalse", "error", err)
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer cursor.Close(ctx)

	var results []yanolja.Product
	if err := cursor.All(ctx, &results); err != nil {
		level.Error(r.logger).Log("repository", "FetchProductsWithContentScheduleStatusFalse", "decode_error", err)
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	// Handle "no documents found" case
	if len(results) == 0 {
		level.Warn(r.logger).Log("repository", "FetchProductsWithContentScheduleStatusFalse", "warning", "no matching documents")
		return nil, mongo.ErrNoDocuments
	}

	return results, nil
}

// GetProductIDVersionAndSalePeriod from product table
func (r *mongoRepository) GetProductIDVersionAndSalePeriod(ctx context.Context) (data []common.ProductValidityAndVersion, err error) {

	level.Info(r.logger).Log(
		"operation", "GetProductIDVersionAndSalePeriod",
	)

	collection := r.db.Collection("products")

	projection := bson.M{
		"productId":      1,
		"productVersion": 1,
		"salePeriod":     1,
		"productName":    1,
		"_id":            0,
	}

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("query error for productId, productVersion, salePeriod, productName: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item common.ProductValidityAndVersion
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("decode error for product: %w", err)
		}
		data = append(data, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return data, nil
}

// ProductNameStruct is a lightweight struct just to decode productName
type ProductNameStruct struct {
	ProductName string `bson:"productName"`
}

func (r *mongoRepository) GetProductNameByProductID(ctx context.Context, productID int64) (string, error) {
	level.Info(r.logger).Log("repository", "GetProductNameByProductID", "product_id", productID)

	// Apply timeout to DB call
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := r.db.Collection("products")

	// Query filter
	filter := bson.M{"productId": productID}

	// Project only the productName field for efficiency
	projection := bson.M{"productName": 1, "_id": 0}

	var result ProductNameStruct
	err := collection.FindOne(ctx, filter,
		options.FindOne().SetProjection(projection),
	).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Warn(r.logger).Log("msg", "no product found", "product_id", productID)
			return "", nil
		}
		level.Error(r.logger).Log("error", "failed to fetch product name", "details", err.Error())
		return "", fmt.Errorf("failed to fetch product name: %w", err)
	}

	return result.ProductName, nil
}
