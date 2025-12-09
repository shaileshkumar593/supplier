package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	domain "swallow-supplier/mongo/domain/travolution"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/utils/constant"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertRequest inserts TicketRequest / BookingRequest / PASSOrPKGRequest into Mongo
func (r *mongoRepository) InsertRequest(ctx context.Context, payload interface{}, requestType string, status string) (string, error) {
	collection := r.db.Collection("travolution_orders")

	now := time.Now().UTC().Format(time.RFC3339)

	order := domain.Order{
		ID:                primitive.NewObjectID().Hex(),
		Status:            status, // default new order status
		CreatedAt:         now,
		UpdatedAt:         now,
		VoucherInfo:       []domain.VoucherInfo{},
		UnitAmounts:       []domain.UnitAmount{},
		BookingStatus:     "",
		BookingAt:         "",
		ExpiredAt:         "",
		ApprovedAt:        "",
		CancelRequestedAt: "",
		CanceledAt:        "",
	}

	switch req := payload.(type) {
	case travolution.TicketRequest:
		order.Type = requestType
		order.Product = req.Product
		order.Option = strconv.Itoa(req.Option)
		order.UnitAmounts = req.UnitAmounts
		order.ReferenceNumber = req.ReferenceNumber
		order.VoucherType = int(req.VoucherSendType)

		order.TravelerName = req.TravelerName
		order.TravelerContactEmail = req.TravelerContactEmail
		order.TravelerContactNumber = req.TravelerContactNumber
		order.TravelerNationality = req.TravelerNationality

	case travolution.BookingRequest:
		order.Type = requestType
		order.Product = req.Product
		order.Option = strconv.Itoa(req.Option)
		order.UnitAmounts = req.UnitAmounts
		order.ReferenceNumber = req.ReferenceNumber
		order.VoucherType = int(req.VoucherSendType)

		order.BookingDate = req.BookingDate
		order.BookingTime = req.BookingTime
		order.BookingAt = strings.TrimSpace(req.BookingDate + " " + req.BookingTime)

		if req.BookingAdditionalInfo != nil {
			raw, err := bson.Marshal(req.BookingAdditionalInfo)
			if err != nil {
				return "", fmt.Errorf("failed to marshal bookingAdditionalInfo: %w", err)
			}
			order.BookingAdditionalInfo = raw
		}

		order.TravelerName = req.TravelerName
		order.TravelerContactEmail = req.TravelerContactEmail
		order.TravelerContactNumber = req.TravelerContactNumber
		order.TravelerNationality = req.TravelerNationality

	case travolution.PASSOrPKGRequest:
		order.Type = requestType // or "PKG" depending on upstream logic
		order.Product = req.Product
		order.Option = req.Option
		order.UnitAmounts = req.UnitAmounts
		order.ReferenceNumber = req.ReferenceNumber
		order.VoucherType = int(req.VoucherSendType)

		order.TravelerName = req.TravelerName
		order.TravelerContactEmail = req.TravelerContactEmail
		order.TravelerContactNumber = req.TravelerContactNumber
		order.TravelerNationality = req.TravelerNationality

	default:
		return "", fmt.Errorf("unsupported payload type %T", req)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := collection.InsertOne(ctx, order); err != nil {
		return "", err
	}

	return order.ID, nil
}

/*

// --- helper: map request_response.UnitAmount -> domain.UnitAmount ---
func convertUnitAmounts(rr []travolution.UnitAmount) []domain.UnitAmount {
	res := make([]domain.UnitAmount, len(rr))
	for i, u := range rr {
		res[i] = domain.UnitAmount{
			Unit:   u.Unit,
			Amount: u.Amount,
		}
	}
	return res
}*/

// for new record createdAt and UpdatedAt must be same

func (r *mongoRepository) UpsertTravolutionOrder(ctx context.Context, payload interface{}, requestType string, status string) (string, error) {
	level.Info(r.logger).Log("repository method", "UpsertTravolutionOrder")

	collection := r.db.Collection("travolution_orders")
	now := time.Now().UTC().Format(time.RFC3339)
	newID := primitive.NewObjectID().Hex() // new Mongo ID

	order := domain.Order{
		ID:                newID,
		Status:            status,
		VoucherInfo:       []domain.VoucherInfo{},
		UnitAmounts:       []domain.UnitAmount{},
		EventHistory:      []domain.WebHookData{},
		BookingStatus:     "",
		BookingAt:         "",
		ExpiredAt:         "",
		ApprovedAt:        "",
		CancelRequestedAt: "",
		CanceledAt:        "",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	req, ok := payload.(travolution.OrderRequest)
	if !ok {
		return "", fmt.Errorf("invalid payload type, expected OrderRequest")
	}

	switch requestType {
	case "TK", "PAS", "PKG":
		order.Type = requestType
		fillOrderCommonFields(&order, req)

	case "BK":
		order.Type = "BK"
		fillOrderCommonFields(&order, req)

		order.BookingDate = req.BookingDate
		order.BookingTime = req.BookingTime
		order.BookingAt = strings.TrimSpace(req.BookingDate + " " + req.BookingTime)
		order.BookingStatus = constant.BookingStatusPending

		// Correctly marshal BookingAdditionalInfo
		if req.BookingAdditionalInfo != nil {
			order.BookingAdditionalInfo = req.BookingAdditionalInfo
		} else {
			order.BookingAdditionalInfo = nil
		}

	default:
		return "", fmt.Errorf("unsupported requestType: %s", requestType)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"referenceNumber": order.ReferenceNumber,
		"orderNumber":     order.OrderNumber, // every time order number is unique
	}

	update := bson.M{
		"$set": bson.M{
			"type":                  order.Type,
			"product":               order.Product,
			"option":                order.Option,
			"referenceNumber":       order.ReferenceNumber,
			"orderNumber":           order.OrderNumber,
			"unitAmounts":           order.UnitAmounts,
			"voucherType":           order.VoucherType,
			"travelerName":          order.TravelerName,
			"travelerContactEmail":  order.TravelerContactEmail,
			"travelerContactNumber": order.TravelerContactNumber,
			"travelerNationality":   order.TravelerNationality,
			"bookingDate":           order.BookingDate,
			"bookingTime":           order.BookingTime,
			"bookingAt":             order.BookingAt,
			"bookingStatus":         order.BookingStatus,
			"expiredAt":             order.ExpiredAt,
			"approvedAt":            order.ApprovedAt,
			"cancelRequestedAt":     order.CancelRequestedAt,
			"canceledAt":            order.CanceledAt,
			"bookingAdditionalInfo": order.BookingAdditionalInfo,
			"updatedAt":             now,
		},
		"$setOnInsert": bson.M{
			"_id":       newID,
			"createdAt": now,
		},
	}

	opts := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upsert order: %w", err)
	}

	if result.UpsertedID != nil {
		if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
			return oid.Hex(), nil
		}
		return newID, nil
	}

	return order.ID, nil
}

// ------------------------- Helper -------------------------
func fillOrderCommonFields(order *domain.Order, req travolution.OrderRequest) {
	order.Product = req.Product
	order.Option = req.Option
	order.UnitAmounts = req.UnitAmounts
	order.ReferenceNumber = req.ReferenceNumber
	order.VoucherType = int(req.VoucherSendType)
	order.TravelerName = req.TravelerName
	order.TravelerContactEmail = req.TravelerContactEmail
	order.TravelerContactNumber = req.TravelerContactNumber
	order.TravelerNationality = req.TravelerNationality
}

func (r *mongoRepository) UpdateTravolutionOrderById(ctx context.Context, update map[string]any) (id string, err error) {
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

	collection := r.db.Collection("travolution_orders")

	filter := bson.M{"_id": update["_id"].(string)}
	opts := options.Update().SetUpsert(false) // Set upsert to true if you want to insert a new document if no match is found

	result, err := collection.UpdateOne(context.TODO(), filter, updateBson, opts)
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("failed to update document: %w", err)
	}

	val := update["_id"]

	return val.(string), nil
}

// GetOrderByOrderNumber fetches an order by orderNumber from travolution_order
func (r *mongoRepository) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (order domain.Order, err error) {
	logger := log.With(r.logger, "repository", "GetOrderByOrderNumber")

	collection := r.db.Collection("travolution_orders")

	// Apply timeout for DB operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"orderNumber": orderNumber}

	err = collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Warn(logger).Log("msg", "order not found", "orderNumber", orderNumber)
			return domain.Order{}, mongo.ErrNoDocuments
		}
		level.Error(logger).Log("msg", "failed to fetch order", "orderNumber", orderNumber, "err", err)
		return domain.Order{}, err
	}

	level.Info(logger).Log("msg", "order fetched successfully", "orderNumber", orderNumber)
	return order, nil
}

// UpdateOrderStatusOnWebhook
func (r *mongoRepository) UpdateOrderStatusOnWebhook(ctx context.Context, orderID string, status, eventType string) error {
	level.Info(r.logger).Log("repository method", "UpdateOrderStatus", "orderID", orderID, "status", status, "eventType", eventType)

	collection := r.db.Collection("travolution_orders")

	// Determine which date field to update
	var dateField string
	switch status {
	case "AV":
		dateField = "bookingAt"
	case "EP":
		dateField = "expiredAt"
	case "AP":
		dateField = "approvedAt"
	case "CR":
		dateField = "cancelRequestedAt"
	case "CL":
		dateField = "canceledAt"
	case "PC":
		// optional: only update UpdatedAt
		dateField = ""
	default:
		return fmt.Errorf("unsupported status: %s", status)
	}

	// Prepare update document
	updateFields := bson.M{
		"status":    status,
		"eventType": eventType,
		"updatedAt": time.Now().UTC().Format(time.RFC3339),
	}

	if dateField != "" {
		updateFields[dateField] = time.Now().UTC().Format(time.RFC3339)
	}

	update := bson.M{
		"$set": updateFields,
	}

	// Perform update
	res, err := collection.UpdateOne(ctx, bson.M{"_id": orderID}, update)
	if err != nil {
		level.Error(r.logger).Log("repository error", "UpdateOrderStatus failed", "err", err)
		return err
	}

	if res.MatchedCount == 0 {
		level.Error(r.logger).Log("repository error", "no order found", "orderID", orderID)
		return mongo.ErrNoDocuments
	}

	level.Info(r.logger).Log("repository success", "UpdateOrderStatus", "orderID", orderID, "status", status)
	return nil
}

// UpdateOrderByOrderNumber
func (r *mongoRepository) UpdateOrderByOrderNumber(ctx context.Context, orderNumber string, update map[string]any) error {
	level.Info(r.logger).Log(
		"repository method", "UpdateCancelRequestedAt",
		"orderNumber", orderNumber,
	)

	collection := r.db.Collection("travolution_orders")

	// Build $set object
	updateBson := bson.M{"$set": bson.M{}}
	for key, val := range update {
		if key != "_id" { // never update _id
			updateBson["$set"].(bson.M)[key] = val
		}
	}

	// Use orderNumber as filter (not _id)
	filter := bson.M{"orderNumber": orderNumber}

	result, err := collection.UpdateOne(ctx, filter, updateBson)
	if err != nil {
		level.Error(r.logger).Log("error while updating order %w", err)
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		level.Error(r.logger).Log("msg", "no document matched", "orderNumber", orderNumber)
		return fmt.Errorf("no document found with orderNumber %s", orderNumber)
	}

	return nil
}
