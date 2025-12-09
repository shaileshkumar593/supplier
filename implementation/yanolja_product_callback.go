package implementation

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertProductClbk
func (s *service) InsertProductClbk(ctx context.Context, product yanolja.Upsert_Product) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "InsertProducts",
		"Request ID", requestID,
	)

	// Defer for panic recovery
	defer func(context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode")
			resp.Code = "500"
		}
	}(ctx)

	record, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, product.Body.ProductID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no record exists for productId", product.Body.ProductID)
			product.Body.OodoSyncStatus = false
			err = s.mongoRepository[config.Instance().MongoDBName].UpsertProduct(ctx, product.Body)
			if err != nil {
				resp.Code = "500"
				resp.Body = "Error inserting product into the database"
				return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on inserting product into the database, %v", err), "UpsertProduct")

			}

		} else {
			level.Error(logger).Log("repository error", "error in fetching product by ", product.Body.ProductID)
			resp.Code = "500"
			resp.Body = fmt.Sprintf("error in fetching product by productId: %d", product.Body.ProductID)
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching product from the database, %v", err), "FetchProductByProductId")
		}

	}

	// Insert the translated product into the database
	if record.ProductVersion <= product.Body.ProductVersion {
		product.Body.OodoSyncStatus = false
		err = s.mongoRepository[config.Instance().MongoDBName].UpsertProduct(ctx, product.Body)
		if err != nil {
			resp.Code = "500"
			resp.Body = "Error inserting product into the database"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on inserting product into the database, %v", err), "UpsertProduct")

		}
	} else {
		resp.Body = product.Body
		resp.Code = "400"
		return resp, customError.NewError(ctx, "leisure-api-0026", fmt.Sprintf("requested product version in not greater than existing product version , %v", err), nil)
	}

	// Fetch the newly inserted product from the database
	record, err = s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, product.Body.ProductID)
	if err != nil {
		level.Error(logger).Log("error", "no record exists for productId", product.Body.ProductID)
		resp.Code = "500"
		resp.Body = fmt.Sprintf("No product found with productId: %d", product.Body.ProductID)
		return resp, fmt.Errorf("product not found: %w", err)
	}

	level.Info(logger).Log("info", "response from api", record)

	// Return the successful response
	resp.Code = "200"
	resp.Body = record
	return resp, nil
}
