package repository

import (
	"context"
	"fmt"
	"swallow-supplier/mongo/domain/odoo"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertOrder inserts or updates an odoo order data in MongoDB
func (r *mongoRepository) UpsertOdooOrder(ctx context.Context, orders []odoo.Order) ([]odoo.Order, error) {
	level.Info(r.logger).Log("method", "UpsertOdooOrder")

	collection := r.db.Collection("odoo_orders")
	var bulkOps []mongo.WriteModel
	var orderIDs []int64
	//now := time.Now().UTC().Format(time.RFC3339)

	if len(orders) == 0 {
		return nil, fmt.Errorf("order is not available for upsert")
	}

	for _, order := range orders {
		filter := bson.M{"orderId": order.OrderID}

		// Marshal struct to BSON
		orderDoc, err := bson.Marshal(order)
		if err != nil {
			level.Error(r.logger).Log("error", "marshal order failed", "orderId", order.OrderID, "err", err)
			continue
		}

		var orderBson bson.M
		if err := bson.Unmarshal(orderDoc, &orderBson); err != nil {
			level.Error(r.logger).Log("error", "unmarshal order BSON failed", "orderId", order.OrderID, "err", err)
			continue
		}

		// Remove conflicting fields from $set block
		//delete(orderBson, "createdAt")
		//delete(orderBson, "updatedAt")
		delete(orderBson, "_id")

		// Add updatedAt
		//orderBson["updatedAt"] = now

		update := bson.M{
			"$set": orderBson,
			"$setOnInsert": bson.M{
				"_id": primitive.NewObjectID(),
			},
		}

		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		bulkOps = append(bulkOps, model)
		orderIDs = append(orderIDs, order.OrderID)
	}

	if len(bulkOps) == 0 {
		return nil, fmt.Errorf("no valid orders to upsert")
	}

	_, err := collection.BulkWrite(ctx, bulkOps, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return nil, fmt.Errorf("bulk upsert failed: %w", err)
	}

	// Fetch upserted/updated documents
	cursor, err := collection.Find(ctx, bson.M{"orderId": bson.M{"$in": orderIDs}})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated orders: %w", err)
	}
	defer cursor.Close(ctx)

	var updatedOrders []odoo.Order
	for cursor.Next(ctx) {
		var order odoo.Order
		if err := cursor.Decode(&order); err != nil {
			level.Error(r.logger).Log("error", "decode order failed", "err", err)
			continue
		}
		updatedOrders = append(updatedOrders, order)
	}

	if err := cursor.Err(); err != nil {
		level.Error(r.logger).Log("error ", "cursor error")
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	level.Info(r.logger).Log("info", "bulk upsert successful", "orderCount", len(updatedOrders))
	return updatedOrders, nil
}

// odoo product update
func (r *mongoRepository) UpsertOdooProduct(ctx context.Context, products []odoo.Product) ([]odoo.Product, error) {
	level.Info(r.logger).Log("method", "UpsertOdooProduct")

	collection := r.db.Collection("odoo_product")
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)

	var productIDs []interface{}

	for _, product := range products {
		filter := bson.M{"product_id": product.ProductID}

		productDoc, err := bson.Marshal(product)
		if err != nil {
			level.Error(r.logger).Log("error", "marshal product failed", "productId", product.ProductID, "err", err)
			continue
		}

		var productBson bson.M
		if err := bson.Unmarshal(productDoc, &productBson); err != nil {
			level.Error(r.logger).Log("error", "unmarshal bson failed", "productId", product.ProductID, "err", err)
			continue
		}

		// Clean system-managed fields
		delete(productBson, "createdAt")
		delete(productBson, "updatedAt")
		delete(productBson, "_id")
		delete(productBson, "product_id")

		productBson["updatedAt"] = nowStr

		update := bson.M{
			"$set": productBson,
			"$setOnInsert": bson.M{
				"_id":        primitive.NewObjectID(),
				"product_id": product.ProductID,
				"createdAt":  nowStr,
			},
		}

		opts := options.Update().SetUpsert(true)

		if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
			level.Error(r.logger).Log("error", "product upsert failed", "productId", product.ProductID, "err", err)
			return nil, fmt.Errorf("failed to upsert product %d: %w", product.ProductID, err)
		}

		productIDs = append(productIDs, product.ProductID)
	}

	// Fetch the updated documents
	cursor, err := collection.Find(ctx, bson.M{"product_id": bson.M{"$in": productIDs}})
	if err != nil {
		level.Error(r.logger).Log("error", "failed to retrieve updated products", "err", err)
		return nil, fmt.Errorf("failed to fetch updated products: %w", err)
	}
	defer cursor.Close(ctx)

	var updatedProducts []odoo.Product
	if err := cursor.All(ctx, &updatedProducts); err != nil {
		level.Error(r.logger).Log("error", "failed to decode updated products", "err", err)
		return nil, fmt.Errorf("failed to decode updated products: %w", err)
	}

	level.Info(r.logger).Log("info", "product upsert completed", "count", len(updatedProducts))
	return updatedProducts, nil
}

// GetOdooOrderbyOrderId  get order from odoo_orders
func (r *mongoRepository) GetOdooOrderbyOrderId(ctx context.Context, orderid int64) (record odoo.Order, err error) {
	level.Info(r.logger).Log("repository method ", "GetOdooOrderbyOrderId")

	collection := r.db.Collection("odoo_orders")

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

// UpdateOdooSyncStatusToFalseIfTrue
func (r *mongoRepository) UpdateOdooSyncStatusToFalseIfTrue(ctx context.Context) error {
	const method = "UpdateOdooSyncStatusToFalseIfTrue"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := r.db.Collection("orders")

	filter := bson.M{"oodoSyncStatus": true}
	update := bson.M{"$set": bson.M{"oodoSyncStatus": false}}

	opts := options.Update().SetUpsert(false)

	result, err := collection.UpdateMany(ctx, filter, update, opts)
	if err != nil {
		level.Error(r.logger).Log("error", err)
		return err
	}

	log.With(r.logger, "method", method).Log(
		"matchedCount", result.MatchedCount,
		"modifiedCount", result.ModifiedCount,
	)

	return nil
}
