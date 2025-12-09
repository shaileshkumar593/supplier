package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	domain "swallow-supplier/mongo/domain/travolution"
	"swallow-supplier/request_response/travolution"
	"time"

	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* func (r *mongoRepository) InsertTravolutionProduct(ctx context.Context, rawProduct travolution.RawProduct) (id string, err error) {
	level.Info(r.logger).Log("repository method", "InsertTravolutionProduct")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convert RawProduct to Product
	product := domain.Product{
		ProductUID:               rawProduct.UID,
		SupplierName:             "Travolution",
		Type:                     rawProduct.Type,
		Status:                   rawProduct.Status,
		SaleTarget:               rawProduct.SaleTarget,
		HasBookingAdditionalInfo: rawProduct.HasBookingAdditionalInfo,
		VoucherType:              rawProduct.VoucherType,
		Titles:                   rawProduct.Titles,
		Images:                   rawProduct.Images,
		Contents:                 rawProduct.Contents,
		Options:                  convertOptions(rawProduct.Options),
		CreatedAt:                time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:                time.Now().UTC().Format(time.RFC3339),
	}

	collection := r.db.Collection("travolution_products")

	result, err := collection.InsertOne(ctx, product, options.InsertOne())
	if err != nil {
		level.Error(r.logger).Log("error", "failed to insert product", "details", err.Error())
		return "", err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		level.Error(r.logger).Log("error", "inserted ID is not ObjectID")
		return "", errors.New("inserted ID is not an ObjectID")
	}

	level.Info(r.logger).Log("inserted_id", insertedID.Hex())
	return insertedID.Hex(), nil
} */

func (r *mongoRepository) UpsertTravolutionProduct(ctx context.Context, rawProduct travolution.RawProduct) (string, error) {
	level.Info(r.logger).Log("repository", "UpsertTravolutionProduct")

	// Apply timeout for DB operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := r.db.Collection("travolution_products")
	now := time.Now().UTC().Format(time.RFC3339)

	// Prepare full product document (all overwriteable fields)
	product := bson.M{
		"productUid":               rawProduct.UID,
		"supplierName":             "Travolution",
		"type":                     rawProduct.Type,
		"status":                   rawProduct.Status,
		"saleTarget":               rawProduct.SaleTarget,
		"hasBookingAdditionalInfo": rawProduct.HasBookingAdditionalInfo,
		"voucherType":              rawProduct.VoucherType,
		"titles":                   rawProduct.Titles,
		"images":                   rawProduct.Images,
		"contents":                 rawProduct.Contents,
		"options":                  convertOptions(rawProduct.Options),
		"updatedAt":                now,
	}

	// MongoDB update document
	update := bson.M{
		"$set":         product,                  // overwrite all fields
		"$setOnInsert": bson.M{"createdAt": now}, // only set CreatedAt if inserting
	}

	// Upsert options
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	// Perform upsert
	var updatedDoc domain.Product
	err := collection.FindOneAndUpdate(ctx,
		bson.M{"productUid": rawProduct.UID}, // filter by UID
		update,
		opts,
	).Decode(&updatedDoc)

	if err != nil {
		level.Error(r.logger).Log(
			"method", "UpsertTravolutionProduct",
			"product_uid", rawProduct.UID,
			"error", err.Error(),
		)
		return "", err
	}

	level.Info(r.logger).Log(
		"method", "UpsertTravolutionProduct",
		"product_id", updatedDoc.Id,
		"product_uid", rawProduct.UID,
		"status", "inserted/updated",
	)

	return updatedDoc.Id, nil
}

func convertOptions(rawOptions []travolution.RawProductOption) []domain.Option {
	options := make([]domain.Option, 0, len(rawOptions))
	for _, ro := range rawOptions {
		uid, _ := handleUID(ro.UID)
		option := domain.Option{
			OptionUID:               uid,
			Names:                   ro.Names,
			Notice:                  ro.Notice,
			UnitAndPriceDetails:     convertUnitPrices(ro.UnitsPrice),
			BookingSchedules:        ro.BookingSchedules,
			AdditionalBookingDetail: convertAdditionalInfo(ro.AdditionalBookingInfo),
		}
		options = append(options, option)
	}
	return options
}

func convertUnitPrices(units []travolution.RawOptionUnitPrice) []domain.UnitAndPriceDetail {
	details := make([]domain.UnitAndPriceDetail, 0, len(units))
	for _, u := range units {
		uid, _ := handleUID(u.UID)
		details = append(details, domain.UnitAndPriceDetail{
			UnitUID:       uid,
			Currency:      u.Currency,
			OriginalPrice: float32(u.OriginalPrice),
			B2BPrice:      float32(u.B2BPrice),
			B2CPrice:      float32(u.B2CPrice),
			MinAmount:     float32(u.MinAmount),
			MaxAmount:     float32(u.MaxAmount),
			Names:         u.Names,
		})
	}
	return details
}

func convertAdditionalInfo(rawInfo []travolution.RawBookingAdditionalInfo) []domain.BookingAdditionalInfo {
	infoList := make([]domain.BookingAdditionalInfo, 0, len(rawInfo))
	for _, info := range rawInfo {
		infoList = append(infoList, domain.BookingAdditionalInfo{
			AdditionalInfoUID: info.UID,
			Type:              info.Type,
			AnswerType:        info.AnswerType,
			Titles:            info.Titles,
			Options:           info.Options,
		})
	}
	return infoList
}

func handleUID(uid interface{}) (string, error) {
	switch v := uid.(type) {
	case string: // already a string
		return v, nil
	case json.Number: // if you used UseNumber()
		return v.String(), nil
	case float64: // default json decoding for numbers
		return strconv.Itoa(int(v)), nil
	case int: // rarely, if unmarshalled directly
		return strconv.Itoa(v), nil
	default:
		return "", fmt.Errorf("unsupported UID type: %T", uid)
	}
}

// GetByProductUid fetches a product by its productUid.
func (r *mongoRepository) GetProductByProductUid(ctx context.Context, productUid int) (product domain.Product, err error) {
	level.Info(r.logger).Log("repository", "GetProductByProductUid")

	collection := r.db.Collection("travolution_products")

	filter := bson.M{"productUid": productUid}
	err = collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Product{}, mongo.ErrNoDocuments // Not found
		}
		return domain.Product{}, fmt.Errorf("failed to fetch productUid=%d: %w", productUid, err)
	}

	return product, nil
}
